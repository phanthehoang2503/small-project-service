package broker

import (
	"context"
	"errors"
)

func (b *Broker) PublishOrderRequested(ctx context.Context, exchange, routingKey string, payload any) error {
	return b.PublishJSON(ctx, exchange, routingKey, payload)
}

func PublishOrderRequested(ctx context.Context, exchange, routingKey string, payload any) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.PublishJSON(ctx, exchange, routingKey, payload)
}
