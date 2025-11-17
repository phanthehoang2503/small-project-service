package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/logger-service/internal/handler"
	"github.com/phanthehoang2503/small-project/logger-service/internal/logger"
	"github.com/rs/zerolog"
)

func main() {
	zlog := logger.InitLogger("logger-service", "./logs/app.log", zerolog.DebugLevel)

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	go func() {
		exchange := event.ExchangeLogs
		queue := "logger_queue"
		bindKey := "#"

		if err := logger.StartConsumer(amqpURL, exchange, queue, bindKey); err != nil {
			zlog.Error().Err(err).Msg("Failed to start RabbitMQ consumer")
		}
	}()

	r := gin.Default()
	h := handler.NewHandler(zlog)
	r.POST("/ingest", h.ReceiveLog)

	zlog.Info().Msg("logger-service running on port 8085")
	if err := r.Run(":8085"); err != nil {
		zlog.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}
