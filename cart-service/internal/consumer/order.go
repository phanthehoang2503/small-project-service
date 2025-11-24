package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/internal/broker"
)

type OrderConsumer struct {
	cartRepo *repo.CartRepo
	broker   *broker.Broker
}

func NewOrderConsumer(cr *repo.CartRepo, b *broker.Broker) *OrderConsumer {
	return &OrderConsumer{
		cartRepo: cr,
		broker:   b,
	}
}

func (c *OrderConsumer) Start(exchange, queueName, routingKey string) error {
	if err := c.broker.DeclareTopicExchange(exchange); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if err := c.broker.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := c.broker.BindQueue(queueName, exchange, []string{routingKey}); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	if err := c.broker.Consume(queueName, c.handleOrderRequested); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("OrderConsumer started, listening on %s -> %s", exchange, queueName)
	return nil
}

type orderRequestedPayload struct {
	CorrelationID string `json:"correlation_id"`
	OrderUUID     string `json:"order_uuid"`
	UserID        uint   `json:"user_id"`
	Total         int64  `json:"total"`
	Currency      string `json:"currency"`
}

func (c *OrderConsumer) handleOrderRequested(routingKey string, body []byte) error {
	var payload orderRequestedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("Received order.requested for user %d (order %s)", payload.UserID, payload.OrderUUID)

	if err := c.cartRepo.ClearCart(payload.UserID); err != nil {
		return fmt.Errorf("failed to clear cart for user %d: %w", payload.UserID, err)
	}

	log.Printf("Cart cleared for user %d", payload.UserID)
	return nil
}
