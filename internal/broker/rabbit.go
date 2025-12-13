package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type Broker struct {
	url          string
	conn         *amqp.Connection
	mu           sync.Mutex
	pubChan      *amqp.Channel
	pubChanMutex sync.Mutex
}

var Global *Broker

// Init initializes the broker and establishes the first connection.
func Init(url string) (*Broker, error) {
	if url == "" {
		return nil, errors.New("rabbitmq url required")
	}

	b := &Broker{
		url: url,
	}

	if err := b.connect(); err != nil {
		return nil, err
	}

	if err := b.SetupDLX(); err != nil {
		return nil, err
	}

	// Start background reconnection handler
	go b.watchConnection()

	Global = b
	return b, nil
}

// connect establishes the connection and sets up the publishing channel.
func (b *Broker) connect() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	log.Printf("RabbitMQ: attempting to connect to %s", b.url)
	conn, err := amqp.Dial(b.url)
	if err != nil {
		return err
	}
	b.conn = conn

	// Setup publishing channel
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	b.pubChan = ch

	log.Println("RabbitMQ: connected")
	return nil
}

// watchConnection monitors the connection and reconnects on failure.
func (b *Broker) watchConnection() {
	for {
		b.mu.Lock()
		if b.conn == nil {
			b.mu.Unlock()
			time.Sleep(2 * time.Second)
			continue
		}
		notifyClose := b.conn.NotifyClose(make(chan *amqp.Error))
		b.mu.Unlock()

		// Block until close notification
		err := <-notifyClose
		if err != nil {
			log.Printf("RabbitMQ: connection lost: %v", err)
			b.reconnect()
		} else {
			// Graceful shutdown
			log.Println("RabbitMQ: connection closed gracefully")
			return
		}
	}
}

func (b *Broker) reconnect() {
	for {
		time.Sleep(3 * time.Second)
		if err := b.connect(); err != nil {
			log.Printf("RabbitMQ: reconnection failed: %v", err)
			continue
		}
		log.Println("RabbitMQ: reconnected")
		return
	}
}

// getPubChannel returns the current publishing channel safely.
func (b *Broker) getPubChannel() (*amqp.Channel, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.conn == nil || b.conn.IsClosed() {
		return nil, errors.New("connection closed")
	}
	if b.pubChan == nil || b.pubChan.IsClosed() {
		ch, err := b.conn.Channel()
		if err != nil {
			return nil, err
		}
		b.pubChan = ch
	}
	return b.pubChan, nil
}

// DeclareTopicExchange declares a durable topic exchange.
func (b *Broker) DeclareTopicExchange(name string) error {
	ch, err := b.getPubChannel()
	if err != nil {
		return err
	}
	b.pubChanMutex.Lock()
	defer b.pubChanMutex.Unlock()
	return ch.ExchangeDeclare(name, "topic", true, false, false, false, nil)
}

// DeclareQueue declares a durable queue.
func (b *Broker) DeclareQueue(name string) error {
	ch, err := b.getPubChannel()
	if err != nil {
		return err
	}
	b.pubChanMutex.Lock()
	defer b.pubChanMutex.Unlock()
	// Arguments for DLQ
	args := amqp.Table{
		"x-dead-letter-exchange":    "dlx",
		"x-dead-letter-routing-key": "dlq",
	}
	_, err = ch.QueueDeclare(name, true, false, false, false, args)
	return err
}

// BindQueue binds a queue to an exchange.
func (b *Broker) BindQueue(queue, exchange string, routingKeys []string) error {
	ch, err := b.getPubChannel()
	if err != nil {
		return err
	}
	b.pubChanMutex.Lock()
	defer b.pubChanMutex.Unlock()

	for _, key := range routingKeys {
		if err := ch.QueueBind(queue, key, exchange, false, nil); err != nil {
			return err
		}
	}
	return nil
}

// PublishJSON publishes a JSON payload with Context Propagation.
func (b *Broker) PublishJSON(ctx context.Context, exchange, routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Inject Trace Context
	headers := make(amqp.Table)
	otel.GetTextMapPropagator().Inject(ctx, AMQPCarrier(headers))

	// Retry logic for publishing
	for i := 0; i < 3; i++ {
		ch, err := b.getPubChannel()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		b.pubChanMutex.Lock()
		err = ch.Publish(
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
				Headers:     headers,
				Timestamp:   time.Now(),
			},
		)
		b.pubChanMutex.Unlock()

		if err == nil {
			return nil
		}
		log.Printf("RabbitMQ: publish failed (attempt %d): %v", i+1, err)
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("failed to publish message after retries")
}

// Consume starts consuming messages using a NEW dedicated channel.
func (b *Broker) Consume(queue string, handler func(ctx context.Context, routingKey string, body []byte) error) error {
	b.mu.Lock()
	if b.conn == nil || b.conn.IsClosed() {
		b.mu.Unlock()
		return errors.New("connection closed")
	}
	// Create a NEW channel for this consumer
	ch, err := b.conn.Channel()
	b.mu.Unlock()

	if err != nil {
		return err
	}

	if err := ch.Qos(10, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("RabbitMQ: consuming from queue %q", queue)
		for msg := range msgs {
			// Extract Trace Context
			if msg.Headers == nil {
				msg.Headers = make(amqp.Table)
			}
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), AMQPCarrier(msg.Headers))

			if err := handler(ctx, msg.RoutingKey, msg.Body); err != nil {
				// FAILURE STRATEGY: Do NOT requeue. Send to DLX (if configured).
				// If no DLX, message is dropped.
				_ = msg.Nack(false, false)
				continue
			}
			_ = msg.Ack(false)
		}
		log.Printf("RabbitMQ: consumer for queue %q stopped", queue)
	}()

	return nil
}

// SetupDLX configures the Dead Letter Exchange and Queue.
func (b *Broker) SetupDLX() error {
	ch, err := b.getPubChannel()
	if err != nil {
		return err
	}

	// 1. Declare DLX
	if err := ch.ExchangeDeclare("dlx", "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// 2. Declare DLQ
	if _, err := ch.QueueDeclare("dlq", true, false, false, false, nil); err != nil {
		return err
	}

	// 3. Bind DLQ to DLX
	if err := ch.QueueBind("dlq", "dlq", "dlx", false, nil); err != nil {
		return err
	}

	return nil
}

// Close closes the connection.
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.conn != nil {
		b.conn.Close()
	}
}

// --- Global Helpers ---

func DeclareTopicExchange(exchange string) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.DeclareTopicExchange(exchange)
}

func DeclareQueue(name string) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.DeclareQueue(name)
}

func BindQueue(queue, exchange string, keys []string) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.BindQueue(queue, exchange, keys)
}

func PublishJSON(ctx context.Context, exchange, routingKey string, payload any) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.PublishJSON(ctx, exchange, routingKey, payload)
}

func Consume(queue string, handler func(ctx context.Context, routingKey string, body []byte) error) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.Consume(queue, handler)
}

func Channel() *amqp.Channel {
	if Global == nil {
		return nil
	}
	ch, _ := Global.getPubChannel()
	return ch
}

// AMQPCarrier adapts amqp.Table to propagation.TextMapCarrier
type AMQPCarrier amqp.Table

func (c AMQPCarrier) Get(key string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (c AMQPCarrier) Set(key string, value string) {
	c[key] = value
}

func (c AMQPCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
