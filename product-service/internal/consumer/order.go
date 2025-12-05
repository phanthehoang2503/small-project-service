package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
)

type OrderConsumer struct {
	repo *repo.Database
	b    *broker.Broker
}

func NewOrderConsumer(r *repo.Database, b *broker.Broker) *OrderConsumer {
	return &OrderConsumer{
		repo: r,
		b:    b,
	}
}

func (c *OrderConsumer) Start(queueName string) error {
	return c.b.Consume(queueName, c.handle)
}

func (c *OrderConsumer) handle(ctx context.Context, routingKey string, body []byte) error {
	if routingKey != event.RoutingKeyOrderRequested {
		return nil
	}

	var payload message.OrderRequested
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[product-consumer] failed to unmarshal order.requested: %v", err)
		return nil // ack to avoid loop
	}

	log.Printf("[product-consumer] received order.requested order=%s items=%d", payload.OrderUUID, len(payload.Items))

	// Deduct stock for each item
	for _, item := range payload.Items {
		if err := c.repo.DeductStock(item.ProductID, item.Quantity); err != nil {
			log.Printf("[product-consumer] failed to deduct stock for product %d: %v", item.ProductID, err)

			// Publish stock.failed event to cancel order
			failEvent := message.StockFailed{
				OrderUUID: payload.OrderUUID,
				Reason:    err.Error(),
			}
			if pubErr := c.b.PublishJSON(ctx, event.ExchangeOrder, event.RoutingKeyStockFailed, failEvent); pubErr != nil {
				log.Printf("[product-consumer] failed to publish stock.failed event: %v", pubErr)
			} else {
				log.Printf("[product-consumer] published stock.failed for order %s", payload.OrderUUID)
			}

			return nil // Return nil to ack the message
		}
	}

	log.Printf("[product-consumer] stock deducted for order %s", payload.OrderUUID)
	return nil
}
