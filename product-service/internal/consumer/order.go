package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type OrderConsumer struct {
	repo  *repo.Database
	cache *repo.CacheRepository
	b     *broker.Broker
}

func NewOrderConsumer(r *repo.Database, c *repo.CacheRepository, b *broker.Broker) *OrderConsumer {
	return &OrderConsumer{
		repo:  r,
		cache: c,
		b:     b,
	}
}

func (c *OrderConsumer) Start(queueName string) error {
	return c.b.Consume(queueName, c.handle)
}

func (c *OrderConsumer) handle(ctx context.Context, routingKey string, body []byte) error {
	tr := otel.Tracer("product-service")
	ctx, span := tr.Start(ctx, "consumer.Handle")
	defer span.End()

	// 1. Handle Order Created (Stock Reservation)
	if routingKey == event.RoutingKeyOrderCreated {
		var payload message.OrderRequested
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("[product-consumer] failed to unmarshal order.created: %v", err)
			return nil
		}
		log.Printf("[product-consumer] received order.created order=%s items=%d", payload.OrderUUID, len(payload.Items))

		var stockItems []repo.StockItem
		for _, item := range payload.Items {
			stockItems = append(stockItems, repo.StockItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		// Try to deduct
		if err := c.repo.BatchDeductStock(stockItems); err != nil {
			log.Printf("[product-consumer] failed to deduct stock: %v", err)

			span.RecordError(err)
			span.SetStatus(codes.Error, "stock_deduction_failed")

			// Publish Failed
			failEvent := message.InventoryReservationFailed{
				OrderUUID: payload.OrderUUID,
				Reason:    err.Error(),
			}
			c.b.PublishJSON(ctx, event.ExchangeOrder, event.RoutingKeyInventoryReservationFailed, failEvent)
			return nil
		}

		// Success -> Publish Reserved
		successEvent := message.InventoryReserved{
			CorrelationID: payload.CorrelationID,
			OrderUUID:     payload.OrderUUID,
			UserID:        payload.UserID,
			Total:         payload.Total,
			Currency:      payload.Currency,
		}
		if err := c.b.PublishJSON(ctx, event.ExchangeOrder, event.RoutingKeyInventoryReserved, successEvent); err != nil {
			log.Printf("[product-consumer] failed to publish inventory.reserved: %v", err)
			// TODO: If publish fails, we should rollback stock immediately?
			// For simple demo, we log error. In production, this is critical.
		}

		// Invalidate Cache
		for _, item := range payload.Items {
			if c.cache != nil {
				c.cache.InvalidateProduct(ctx, item.ProductID)
			}
		}
		log.Printf("[product-consumer] stock reserved & event published for order %s", payload.OrderUUID)
		return nil
	}

	// 2. Handle Payment Failed (Compensation / Rollback)
	if routingKey == event.RoutingKeyPaymentFailed {
		// TODO: Implement rollback logic.
		// Currently, the payment.failed event does not contain item details needed to reverse stock.
		log.Printf("[product-consumer] received payment.failed - rollback skipped (missing item details)")
		return nil
	}

	return nil
}
