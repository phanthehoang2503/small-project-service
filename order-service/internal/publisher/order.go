package publisher

import (
	"context"
	"log"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
)

func PublishOrderCreated(ctx context.Context, b *broker.Broker, correlationID, orderUUID string, userID uint, total int64, currency string, items []message.OrderItem) error {
	payload := message.OrderRequested{
		CorrelationID: correlationID,
		OrderUUID:     orderUUID,
		UserID:        userID,
		Total:         total,
		Currency:      currency,
		Items:         items,
	}

	if err := b.PublishJSON(ctx, event.ExchangeOrder, event.RoutingKeyOrderCreated, payload); err != nil {
		log.Printf("[order-publisher] failed to publish order.created: %v", err)
		return err
	}
	log.Printf("[order-publisher] published order.created for order %s", orderUUID)
	return nil
}
