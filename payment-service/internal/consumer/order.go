package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/payment-service/internal/repo"
)

type PaymentConsumer struct {
	repo *repo.PaymentRepo
	b    *broker.Broker
}

// creates consumer
func NewPaymentConsumer(r *repo.PaymentRepo, b *broker.Broker) *PaymentConsumer {
	return &PaymentConsumer{repo: r, b: b}
}

type orderRequestedPayload struct {
	CorrelationID string `json:"correlation_id"`
	OrderUUID     string `json:"order_uuid"`
	UserID        uint   `json:"user_id"`
	Total         int64  `json:"total"`
	Currency      string `json:"currency"`
}

// Start registers consume handler on the broker. It expects the queue to be declared & bound beforehand.
func (pc *PaymentConsumer) Start(queueName string) error {
	return pc.b.Consume(queueName, pc.handle)
}

// handle implements the signature expected by broker.Consume: func(routingKey string, body []byte) error
func (pc *PaymentConsumer) handle(routingKey string, body []byte) error {
	// expect order.requested
	if routingKey != event.RoutingKeyOrderRequested {
		log.Printf("[payment-consumer] unexpected routing key: %s", routingKey)
		return nil
	}

	var payload orderRequestedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("[payment-consumer] invalid payload: %v", err)
		return nil
	}

	log.Printf("[payment-consumer] processing order.requested order=%s amount=%d", payload.OrderUUID, payload.Total)

	if _, err := pc.repo.CreatePending(payload.OrderUUID, payload.Total, payload.Currency); err != nil {
		log.Printf("[payment-consumer] db create failed: %v", err)
		logger.Error(context.Background(), fmt.Sprintf("Payment failed (db create): order=%s err=%v", payload.OrderUUID, err))
		pc.publishFailure(payload.OrderUUID, payload.CorrelationID, "db_create_failed")
		return err
	}

	time.Sleep(150 * time.Millisecond)

	if err := pc.repo.PaymentSucceeded(payload.OrderUUID); err != nil {
		log.Printf("[payment-consumer] update succeeded failed: %v", err)
		pc.publishFailure(payload.OrderUUID, payload.CorrelationID, "update_status_failed")
		return err
	}

	out := map[string]interface{}{
		"correlation_id": payload.CorrelationID,
		"order_uuid":     payload.OrderUUID,
		"status":         "succeeded",
		"amount":         payload.Total,
		"currency":       payload.Currency,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}

	if err := pc.b.PublishJSON(event.ExchangeOrder, event.RoutingKeyOrderPaid, out); err != nil {
		log.Printf("[payment-consumer] publish order.paid failed: %v", err)
		return err
	}
	log.Printf("[payment-consumer] successfully published order.paid event for order=%s", payload.OrderUUID)

	log.Printf("[payment-consumer] payment succeeded order=%s", payload.OrderUUID)
	logger.Info(context.Background(), fmt.Sprintf("Payment succeeded: order=%s amount=%d", payload.OrderUUID, payload.Total))
	return nil
}

func (pc *PaymentConsumer) publishFailure(orderUUID, correlationID, reason string) {
	out := map[string]interface{}{
		"correlation_id": correlationID,
		"order_uuid":     orderUUID,
		"status":         "failed",
		"reason":         reason,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	}
	if err := pc.b.PublishJSON(event.ExchangeOrder, event.RoutingKeyPaymentFailed, out); err != nil {
		log.Printf("[payment-consumer] failed to publish payment.failed: %v", err)
	} else {
		log.Printf("[payment-consumer] published payment.failed order=%s reason=%s", orderUUID, reason)
	}
}
