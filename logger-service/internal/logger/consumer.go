package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/phanthehoang2503/small-project/internal/message"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// consume messages from RabbitMQ and write them using zerolog
func StartConsumer(amqpURL, exchange, queueName, bindingKey string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare exchange
	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Declare durable queue
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind to the exchange with wildcard (listen to all logs)
	if err := ch.QueueBind(
		q.Name,
		bindingKey,
		exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	// Prepare zerolog output to both console & rotating file
	logFile := "./logs/central.log"
	lumber := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     14,
		Compress:   true,
	}
	// Use JSON output for better parsing in Loki/Grafana
	// consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	writer := zerolog.MultiLevelWriter(os.Stdout, lumber)
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(writer).With().Timestamp().Str("service", "logger-service").Logger()

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("Logger Service waiting on queue=%s exchange=%s binding=%s\n", q.Name, exchange, bindingKey)

	for d := range msgs {
		var ev message.LogEvent
		if err := json.Unmarshal(d.Body, &ev); err != nil {
			logger.Error().Err(err).Msg("Failed to decode message")
			continue
		}
		logger.Info().
			Str("source", ev.Service).
			Str("level", ev.Level).
			Msg(ev.Message)
	}
	return nil
}
