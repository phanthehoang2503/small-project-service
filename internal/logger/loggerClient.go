package logger

import (
	"context"
	"time"

	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/message"
)

var (
	service = "unknown-service"
)

// SetConfig lets each service tell the logger which service it is
func SetConfig(_ string, svc string) {
	if svc != "" {
		service = svc
	}
}

// send publishes log entries to RabbitMQ via the broker package
func send(ctx context.Context, level, msg, traceID string, fields map[string]interface{}) {
	ev := message.LogEvent{
		Service: service,
		Level:   level,
		Message: msg,
		Time:    time.Now().Format(time.RFC3339),
		Meta: map[string]interface{}{
			"trace_id": traceID,
			"fields":   fields,
		},
	}

	broker.Publish(service+"."+level, ev)
}

// sugar wrappers
func Info(ctx context.Context, msg, traceID string, fields map[string]interface{}) {
	send(ctx, "info", msg, traceID, fields)
}
func Warn(ctx context.Context, msg, traceID string, fields map[string]interface{}) {
	send(ctx, "warn", msg, traceID, fields)
}
func Error(ctx context.Context, msg, traceID string, fields map[string]interface{}) {
	send(ctx, "error", msg, traceID, fields)
}
