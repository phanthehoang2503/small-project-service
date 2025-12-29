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

	// Init broker (it handles reconnection internally)
	b, err := broker.Init(rabbitURL)
	if err != nil {
		for i := 0; i < 10; i++ {
			log.Printf("attempt %d: connecting to RabbitMQ...", i+1)
			b, err = broker.Init(rabbitURL)
			if err == nil {
				break
			}
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatal("could not connect to RabbitMQ:", err)
		}
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
