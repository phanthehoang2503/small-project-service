package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
)

type StockConsumer struct {
	repo *repo.OrderRepo
	b    *broker.Broker
}

func NewStockConsumer(r *repo.OrderRepo, b *broker.Broker) *StockConsumer {
	return &StockConsumer{
		repo: r,
		b:    b,
	}
}

func (c *StockConsumer) Start(queueName string) error {
	return c.b.Consume(queueName, c.handle)
}

func (c *StockConsumer) handle(ctx context.Context, routingKey string, body []byte) error {
	if routingKey != event.RoutingKeyInventoryReservationFailed {
		return nil
	}

	var payload message.InventoryReservationFailed
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[order-stock-consumer] failed to unmarshal inventory.reservation.failed: %v", err)
		return nil // ack
	}

	log.Printf("[order-stock-consumer] received inventory.reservation.failed for order %s. Reason: %s", payload.OrderUUID, payload.Reason)

	// Cancel the order
	if _, err := c.repo.UpdateStatusByUUID(payload.OrderUUID, "Cancelled"); err != nil {
		log.Printf("[order-stock-consumer] failed to cancel order %s: %v", payload.OrderUUID, err)
		return err // retry
	}

	log.Printf("[order-stock-consumer] order %s cancelled", payload.OrderUUID)
	return nil
}
