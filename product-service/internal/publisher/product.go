package publisher

import (
	"context"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
)

func PublishProductCreated(ctx context.Context, p *model.Product) error {
	msg := message.ProductMessage{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}

	return publishJSON(event.ExchangeProduct, event.RoutingKeyProductCreated, msg)
}

func PublishProductUpdated(ctx context.Context, p *model.Product) error {
	msg := message.ProductMessage{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}

	return publishJSON(event.ExchangeProduct, event.RoutingKeyProductUpdated, msg)
}

func PublishProductDeleted(ctx context.Context, id uint) error {
	msg := message.ProductMessage{
		ID: id,
	}

	return publishJSON(event.ExchangeProduct, event.RoutingKeyProductDeleted, msg)
}

func publishJSON(exchange, rk string, payload any) error {
	return broker.Global.PublishJSON(exchange, rk, payload)
}
