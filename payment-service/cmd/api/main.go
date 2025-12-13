package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/internal/telemetry"

	"github.com/phanthehoang2503/small-project/payment-service/internal/consumer"
	"github.com/phanthehoang2503/small-project/payment-service/internal/model"
	"github.com/phanthehoang2503/small-project/payment-service/internal/repo"
)

func main() {
	godotenv.Load()

	// Init Tracer
	shutdown := telemetry.InitTracer("payment-service")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// migrations
	if err := db.AutoMigrate(&model.Payment{}); err != nil {
		log.Fatalf("Migration failed (payment): %v", err)
	}

	// Rabbit broker
	b := helper.ConnectRabbit()
	defer b.Close()

	// tell logger which service this is
	logger.SetService("payment-service")

	// declare queue & bind it to order exchange routing key
	queueName := "payment_service_queue"
	if err := b.DeclareQueue(queueName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	// bind the queue to order_exchange with routing key order.requested
	if err := b.BindQueue(queueName, event.ExchangeOrder, []string{event.RoutingKeyOrderRequested}); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	// repo + publisher + consumer
	payRepo := repo.NewPaymentRepo(db)
	pc := consumer.NewPaymentConsumer(payRepo, b)

	if err := pc.Start(queueName); err != nil {
		log.Fatalf("failed to start payment consumer: %v", err)
	}

	// small HTTP API for payment lookup
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(otelgin.Middleware("payment-service"))
	r.GET("/payments/:order_uuid", func(c *gin.Context) {
		orderUUID := c.Param("order_uuid")
		p, err := payRepo.GetByOrderUUID(orderUUID)
		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(200, p)
	})

	r.Run(":8086")
}
