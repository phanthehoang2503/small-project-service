package main

import (
	"log"

	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/mailer-service/internal/consumer"
)

func main() {
	// tell logger which service this is
	logger.SetService("mailer-service")

	// Rabbit broker
	b := helper.ConnectRabbit()
	defer b.Close()

	// declare queue
	queueName := "mailer_queue"
	if err := b.DeclareQueue(queueName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	// Bind to order.paid
	if err := b.BindQueue(queueName, event.ExchangeOrder, []string{event.RoutingKeyOrderPaid}); err != nil {
		log.Fatalf("failed to bind queue order.paid: %v", err)
	}

	// Start consumer
	c := consumer.NewMailerConsumer(b)

	log.Println("Mailer service started, waiting for messages...")

	if err := c.Start(queueName); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}

	// Block forever
	select {}
}
