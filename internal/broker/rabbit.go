package broker

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Global broker
var Global *Broker

// Init sets up connection + channel.
func Init(url string) (*Broker, error) {
	if url == "" {
		return nil, errors.New("rabbitmq url required")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	b := &Broker{
		conn:    conn,
		channel: ch,
	}

	Global = b

	log.Println("RabbitMQ connected")
	return b, nil
}

// DeclareTopicExchange declares a durable topic exchange.
func (b *Broker) DeclareTopicExchange(name string) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}
	if name == "" {
		return errors.New("exchange name required")
	}

	return b.channel.ExchangeDeclare(
		name,
		"topic",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
}

// DeclareQueue declares a durable queue.
func (b *Broker) DeclareQueue(name string) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}
	if name == "" {
		return errors.New("queue name required")
	}

	_, err := b.channel.QueueDeclare(
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

// BindQueue binds a queue to an exchange with a list of routing keys.
func (b *Broker) BindQueue(queue, exchange string, routingKeys []string) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}
	if queue == "" {
		return errors.New("queue name required")
	}
	if exchange == "" {
		return errors.New("exchange name required")
	}
	if len(routingKeys) == 0 {
		return errors.New("at least one routing key required")
	}

	for _, key := range routingKeys {
		if err := b.channel.QueueBind(
			queue,
			key,
			exchange,
			false, // noWait
			nil,   // args
		); err != nil {
			return err
		}
	}
	return nil
}

// PublishJSON publishes a JSON payload to an exchange with a routing key.
func (b *Broker) PublishJSON(exchange, routingKey string, payload any) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}
	if exchange == "" {
		return errors.New("exchange name required")
	}
	if routingKey == "" {
		return errors.New("routing key required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return b.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// Consume starts consuming messages from a queue and passes them to handler.
// handler should return error to Nack (and requeue) or nil to Ack.
func (b *Broker) Consume(queue string, handler func(routingKey string, body []byte) error) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}
	if queue == "" {
		return errors.New("queue name required")
	}
	if handler == nil {
		return errors.New("handler required")
	}

	// optional: set basic QoS
	if err := b.channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return err
	}

	msgs, err := b.channel.Consume(
		queue,
		"",    // consumer tag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("RabbitMQ: consuming from queue %q\n", queue)
		for msg := range msgs {
			if err := handler(msg.RoutingKey, msg.Body); err != nil {
				// handler failed → Nack and requeue
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
		log.Printf("RabbitMQ: consumer for queue %q stopped\n", queue)
	}()

	return nil
}

// Global helpers using the Global broker instance.
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

func PublishJSON(exchange, routingKey string, payload any) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.PublishJSON(exchange, routingKey, payload)
}

func Consume(queue string, handler func(routingKey string, body []byte) error) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.Consume(queue, handler)
}

// Channel returns the underlying AMQP channel.
func (b *Broker) Channel() *amqp.Channel {
	if b == nil {
		return nil
	}
	return b.channel
}

// GlobalChannel returns the AMQP channel from the Global broker.
func GlobalChannel() *amqp.Channel {
	if Global == nil {
		return nil
	}
	return Global.channel
}

// Close closes channel and connection.
func (b *Broker) Close() {
	if b == nil {
		return
	}
	if b.channel != nil {
		_ = b.channel.Close()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}
