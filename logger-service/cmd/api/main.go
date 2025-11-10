package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/logger-service/internal/handler"
	"github.com/phanthehoang2503/small-project/logger-service/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func main() {
	zlog := logger.InitLogger("logger-service", "./logs/app.log", zerolog.DebugLevel)

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	// Try connecting with retry
	conn, err := RabbitMQ(amqpURL, 10)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Failed to connect to RabbitMQ after retries")
	}
	defer conn.Close()

	go func() {
		exchange := "logs_exchange"
		queue := "logger_queue"
		bindKey := "#"

		if err := logger.StartConsumer(amqpURL, exchange, queue, bindKey); err != nil {
			zlog.Error().Err(err).Msg("Failed to start RabbitMQ consumer")
		}
	}()

	// Start HTTP API
	r := gin.Default()
	h := handler.NewHandler(zlog)
	r.POST("/ingest", h.ReceiveLog)

	zlog.Info().Msg("logger-service running on port 8085")
	if err := r.Run(":8085"); err != nil {
		zlog.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}

func RabbitMQ(url string, maxAttempts int) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= maxAttempts; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Printf("Connected to RabbitMQ after %d attempt(s)", i)
			return conn, nil
		}
		log.Printf("attempt %d: failed to connect to RabbitMQ (%v)", i, err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to RabbitMQ after %d attempts: %w", maxAttempts, err)
}
