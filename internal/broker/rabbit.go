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

// DeclareTopicExchange is a helper to declare a durable topic exchange.
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
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

func DeclareTopicExchange(exchange string) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.DeclareTopicExchange(exchange)
}

func PublishJSON(exchange, routingKey string, payload any) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.PublishJSON(exchange, routingKey, payload)
}

func (b *Broker) Channel() *amqp.Channel {
	if b == nil {
		return nil
	}
	return b.channel
}

func GlobalChannel() *amqp.Channel {
	if Global == nil {
		return nil
	}
	return Global.channel
}

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
