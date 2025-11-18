package publisher

import (
	"context"
	"encoding/json"

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

	body, _ := json.Marshal(msg)

	return broker.Global.PublishJSON(
		event.ExchangeProduct,
		event.RoutingKeyProductCreated,
		body,
	)
}

func PublishProductUpdated(ctx context.Context, p *model.Product) error {
	msg := message.ProductMessage{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}

	body, _ := json.Marshal(msg)

	return broker.Global.PublishJSON(
		event.ExchangeProduct,
		event.RoutingKeyProductUpdated,
		body,
	)
}

func PublishProductDeleted(ctx context.Context, id uint) error {
	msg := message.ProductMessage{
		ID: id,
	}

	body, _ := json.Marshal(msg)

	return broker.Global.PublishJSON(
		event.ExchangeProduct,
		event.RoutingKeyProductDeleted,
		body,
	)
}
