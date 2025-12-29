package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

// OrderPaidConsumer listens for order.paid and updates order status.
type OrderPaidConsumer struct {
	repo *repo.OrderRepo
	b    *broker.Broker
}

func NewOrderPaidConsumer(r *repo.OrderRepo, b *broker.Broker) *OrderPaidConsumer {
	return &OrderPaidConsumer{repo: r, b: b}
}

type orderPaidPayload struct {
	CorrelationID string `json:"correlation_id"`
	OrderUUID     string `json:"order_uuid"`
	PaymentID     string `json:"payment_id,omitempty"`
	Status        string `json:"status,omitempty"`
	Amount        int64  `json:"amount,omitempty"`
	Currency      string `json:"currency,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
}

// Start registers the consumer on the given queue (queue should be declared & bound beforehand).
func (c *OrderPaidConsumer) Start(queueName string) error {
	return c.b.Consume(queueName, c.handle)
}

func (c *OrderPaidConsumer) handle(ctx context.Context, routingKey string, body []byte) error {
	tr := otel.Tracer("order-service")
	ctx, span := tr.Start(ctx, "consumer.ProcessPaymentEvent")
	defer span.End()

	if routingKey == event.RoutingKeyPaymentFailed {
		return c.handlePaymentFailed(body)
	}

	if routingKey != event.RoutingKeyPaymentSucceeded {
		// ignore unrelated keys but ack
		return nil
	}

	var p orderPaidPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("[payment-event-consumer] invalid payload: %v", err)
		return nil // ack malformed message
	}

	log.Printf("[payment-event-consumer] received payment.succeeded order=%s amount=%d", p.OrderUUID, p.Amount)

	// skip if payment status is not succeeded
	if p.Status != "" && p.Status != "succeeded" && p.Status != "success" {
		log.Printf("[payment-event-consumer] payment status not succeeded, skipping order=%s status=%s", p.OrderUUID, p.Status)
		return nil
	}

	// attempt to mark order as Paid
	ord, err := c.repo.UpdateStatusIfNot(p.OrderUUID, "Paid", "Cancelled")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[payment-event-consumer] order not found (uuid=%s)", p.OrderUUID)
			return nil
		}
		log.Printf("[payment-event-consumer] failed to update order status: %v", err)
		return err
	}
	if ord == nil {
		log.Printf("[payment-event-consumer] skipping Paid update for Cancelled order %s", p.OrderUUID)
		return nil
	}

	log.Printf("[payment-event-consumer] order marked Paid uuid=%s id=%d", ord.UUID, ord.ID)

	_ = p
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (c *OrderPaidConsumer) handlePaymentFailed(body []byte) error {
	var p struct {
		OrderUUID string `json:"order_uuid"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("[payment-event-consumer] invalid failure payload: %v", err)
		return nil
	}

	log.Printf("[payment-event-consumer] received payment.failed order=%s reason=%s", p.OrderUUID, p.Reason)

	if err := c.repo.CompensateOrder(p.OrderUUID, p.Reason); err != nil {
		log.Printf("[payment-event-consumer] failed to compensate order: %v", err)
		return err
	}

	log.Printf("[payment-event-consumer] order cancelled uuid=%s", p.OrderUUID)
	return nil
}
