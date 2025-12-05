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

	return publishJSON(ctx, event.ExchangeProduct, event.RoutingKeyProductCreated, msg)
}

func PublishProductUpdated(ctx context.Context, p *model.Product) error {
	msg := message.ProductMessage{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}

	return publishJSON(ctx, event.ExchangeProduct, event.RoutingKeyProductUpdated, msg)
}

func PublishProductDeleted(ctx context.Context, id uint) error {
	msg := message.ProductMessage{
		ID: id,
	}

	return publishJSON(ctx, event.ExchangeProduct, event.RoutingKeyProductDeleted, msg)
}

func publishJSON(ctx context.Context, exchange, rk string, payload any) error {
	return broker.Global.PublishJSON(ctx, exchange, rk, payload)
}
