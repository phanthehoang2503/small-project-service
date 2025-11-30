package consumer

import (
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

func (c *OrderConsumer) handle(routingKey string, body []byte) error {
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
			// TODO: Publish stock.failed event to cancel order (compensation)
			// For now, we just log it. In a real system, we MUST compensate.
			return err // retry? or fail?
		}
	}

	log.Printf("[product-consumer] stock deducted for order %s", payload.OrderUUID)
	return nil
}
