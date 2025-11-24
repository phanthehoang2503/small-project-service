package broker

import (
	"errors"
)

func (b *Broker) PublishOrderRequested(exchange, routingKey string, payload any) error {
	return b.PublishJSON(exchange, routingKey, payload)
}

func PublishOrderRequested(exchange, routingKey string, payload any) error {
	if Global == nil {
		return errors.New("global broker not initialized")
	}
	return Global.PublishJSON(exchange, routingKey, payload)
}
