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
	ch := broker.Global.Channel()
	if ch == nil {
		return fmt.Errorf("rabbitmq not initialized")
	}

	// declare queue
	_, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(queue, binding, exchange, false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			var ev message.ProductMessage

			if err := json.Unmarshal(m.Body, &ev); err != nil {
				log.Println("cart-service: failed to parse product event:", err)
				continue
			}

			pc.handleProductEvent(m.RoutingKey, ev)
		}
	}()

	log.Println("cart-service: listening for product events...")
	return nil
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
