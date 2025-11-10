package broker

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/phanthehoang2503/small-project/internal/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// Global broker for convenience (so callers can just call broker.Publish(...))
var Global *Broker

// InitRabbit creates a connection/channel and declares the exchange.
// It DOES NOT close the connection/channel — caller should call Close() on the returned Broker when shutting down.
func InitRabbit(url, exchange string) (*Broker, error) {
	if url == "" {
		return nil, errors.New("rabbitmq url required")
	}
	if exchange == "" {
		return nil, errors.New("exchange name required")
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

	// declare exchange using provided name
	if err := ch.ExchangeDeclare(
		exchange, // name (use parameter)
		"topic",  // type -- topic is usually more flexible than fanout
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // args
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	b := &Broker{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}

	// set package-global for convenience
	Global = b

	log.Println("RabbitMQ connected, exchange:", exchange)
	return b, nil
}

// Publish publishes an event to the exchange using the provided routing key.
func (b *Broker) Publish(routingKey string, event message.LogEvent) error {
	if b == nil || b.channel == nil {
		return errors.New("rabbitmq broker not initialized")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.channel.Publish(
		b.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// Package-level wrapper that uses Global broker
func Publish(routingKey string, event message.LogEvent) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.Publish(routingKey, event)
}

// Close closes the channel and connection. Call this at shutdown.
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
