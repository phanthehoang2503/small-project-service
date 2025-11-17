package logger

import (
	"context"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/message"
)

var service = "unknown-service"

func SetService(name string) {
	if name != "" {
		service = name
	}
}

func send(_ context.Context, level, msg string) {
	ev := message.LogEvent{
		Level:   level,
		Service: service,
		Message: msg,
		Time:    time.Now().Format(time.RFC3339),
	}

	var routingKey string
	switch level {
	case "info":
		routingKey = event.RoutingKeyLogInfo
	case "warn":
		routingKey = event.RoutingKeyLogWarn
	case "error":
		routingKey = event.RoutingKeyLogError
	default:
		routingKey = "log." + level
	}

	_ = broker.PublishJSON(event.ExchangeLogs, routingKey, ev)
}

func Info(ctx context.Context, msg string)  { send(ctx, "info", msg) }
func Warn(ctx context.Context, msg string)  { send(ctx, "warn", msg) }
func Error(ctx context.Context, msg string) { send(ctx, "error", msg) }
