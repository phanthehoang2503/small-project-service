package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/message"
)

// ProductConsumer holds dependencies
type ProductConsumer struct {
	repo *repo.ProductRepo
}

// NewProductConsumer creates consumer
func NewProductConsumer(snapshotRepo *repo.ProductRepo) *ProductConsumer {
	return &ProductConsumer{
		repo: snapshotRepo,
	}
}

// Start begins consuming product events
func (pc *ProductConsumer) Start(exchange, queue, binding string) error {
	if err := broker.DeclareQueue(queue); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := broker.BindQueue(queue, exchange, []string{binding}); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return broker.Consume(queue, func(routingKey string, body []byte) error {
		var ev message.ProductMessage
		if err := json.Unmarshal(body, &ev); err != nil {
			log.Println("cart-service: failed to parse product event:", err)
			return nil // ack malformed
		}
		pc.handleProductEvent(routingKey, ev)
		return nil
	})
}

// process events
func (pc *ProductConsumer) handleProductEvent(routingKey string, ev message.ProductMessage) {
	switch routingKey {

	case "product.created", "product.updated":
		snap := model.ProductSnapshot{
			ProductID: ev.ID,
			Name:      ev.Name,
			Price:     ev.Price,
			Stock:     ev.Stock,
		}

		if err := pc.repo.Upsert(snap); err != nil {
			log.Println("cart-service: failed to upsert snapshot:", err)
		} else {
			log.Println("cart-service: snapshot updated for product:", ev.ID)
		}

	case "product.deleted":
		if err := pc.repo.Delete(ev.ID); err != nil {
			log.Println("cart-service: failed to delete snapshot:", err)
		} else {
			log.Println("cart-service: snapshot deleted for product:", ev.ID)
		}
	}
}
