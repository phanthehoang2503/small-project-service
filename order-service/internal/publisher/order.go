package publisher

import (
	"log"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
)

type orderRequestedPayload struct {
	CorrelationID string `json:"correlation_id"`
	OrderUUID     string `json:"order_uuid"`
	UserID        uint   `json:"user_id"`
	Total         int64  `json:"total"`
	Currency      string `json:"currency"`
}

func PublishOrderRequested(b *broker.Broker, correlationID, orderUUID string, userID uint, total int64, currency string) error {
	payload := orderRequestedPayload{
		CorrelationID: correlationID,
		OrderUUID:     orderUUID,
		UserID:        userID,
		Total:         total,
		Currency:      currency,
	}

	if err := b.PublishJSON(event.ExchangeOrder, event.RoutingKeyOrderRequested, payload); err != nil {
		log.Printf("[order-publisher] failed to publish order.requested: %v", err)
		return err
	}
	log.Printf("[order-publisher] published order.requested for order %s", orderUUID)
	return nil
}
