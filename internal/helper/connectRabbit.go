package helper

import (
	"log"
	"os"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
)

func ConnectRabbit() *broker.Broker {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var (
		b   *broker.Broker
		err error
	)

	for i := 0; i < 10; i++ {
		b, err = broker.Init(rabbitURL)
		if err == nil {
			log.Printf("connected to RabbitMQ (attempt %d)\n", i+1)
			break
		}
		log.Printf("attempt %d: failed to connect to RabbitMQ (%v)\n", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("could not connect to RabbitMQ after multiple attempts:", err)
	}

	// Declare exchanges that this service might use.
	// At minimum, declare logs exchange so the shared logger client can publish.
	if err := b.DeclareTopicExchange(event.ExchangeLogs); err != nil {
		log.Fatalf("failed to declare logs exchange: %v", err)
	}
	if err := b.DeclareTopicExchange(event.ExchangeProduct); err != nil {
		log.Fatalf("failed to declare product exchange: %v", err)
	}
	if err := b.DeclareTopicExchange(event.ExchangeOrder); err != nil {
		log.Fatalf("failed to declare order exchange: %v", err)
	}

	log.Println("RabbitMQ ready in service")
	return b
}
