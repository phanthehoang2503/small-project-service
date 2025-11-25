package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
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

func (c *OrderPaidConsumer) handle(routingKey string, body []byte) error {
	if routingKey == event.RoutingKeyPaymentFailed {
		return c.handlePaymentFailed(body)
	}

	if routingKey != event.RoutingKeyOrderPaid {
		// ignore unrelated keys but ack
		return nil
	}

	var p orderPaidPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("[payment-event-consumer] invalid payload: %v", err)
		return nil // ack malformed message
	}

	log.Printf("[payment-event-consumer] received order.paid order=%s amount=%d", p.OrderUUID, p.Amount)

	// skip if payment status is not succeeded
	if p.Status != "" && p.Status != "succeeded" && p.Status != "success" {
		log.Printf("[payment-event-consumer] payment status not succeeded, skipping order=%s status=%s", p.OrderUUID, p.Status)
		return nil
	}

	// attempt to mark order as Paid
	ord, err := c.repo.UpdateStatusByUUID(p.OrderUUID, "Paid")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[payment-event-consumer] order not found (uuid=%s)", p.OrderUUID)
			return nil
		}
		log.Printf("[payment-event-consumer] failed to update order status: %v", err)
		return err
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

	if _, err := c.repo.UpdateStatusByUUID(p.OrderUUID, "Cancelled"); err != nil {
		log.Printf("[payment-event-consumer] failed to cancel order: %v", err)
		return err
	}

	log.Printf("[payment-event-consumer] order cancelled uuid=%s", p.OrderUUID)
	return nil
}
