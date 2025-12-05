package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
)

type MailerConsumer struct {
	b *broker.Broker
}

func NewMailerConsumer(b *broker.Broker) *MailerConsumer {
	return &MailerConsumer{b: b}
}

func (c *MailerConsumer) Start(queueName string) error {
	return c.b.Consume(queueName, c.handle)
}

func (c *MailerConsumer) handle(ctx context.Context, routingKey string, body []byte) error {
	if routingKey == event.RoutingKeyOrderPaid {
		return c.handleOrderPaid(body)
	}
	return nil
}

type orderPaidPayload struct {
	OrderUUID string `json:"order_uuid"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}

func (c *MailerConsumer) handleOrderPaid(body []byte) error {
	var p orderPaidPayload
	if err := json.Unmarshal(body, &p); err != nil {
		log.Printf("[mailer] invalid payload: %v", err)
		return nil
	}

	log.Printf("[mailer] sending confirmation for order %s", p.OrderUUID)

	// Send email via MailHog
	from := "noreply@example.com"
	to := []string{"customer@example.com"} // In real app, this would come from payload or user service lookup
	msg := []byte(fmt.Sprintf("To: customer@example.com\r\n"+
		"Subject: Order Confirmation %s\r\n"+
		"\r\n"+
		"Thank you for your order!\r\n"+
		"Order ID: %s\r\n"+
		"Total: %d %s\r\n", p.OrderUUID, p.OrderUUID, p.Amount, p.Currency))

	// MailHog is available at 'mailhog:1025' inside docker network
	err := smtp.SendMail("mailhog:1025", nil, from, to, msg)
	if err != nil {
		log.Printf("[mailer] failed to send email: %v", err)
		return err
	}

	log.Printf("[mailer] email sent for order %s", p.OrderUUID)
	return nil
}
