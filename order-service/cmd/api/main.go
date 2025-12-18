package main

import (
	"log"
	"os"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/internal/telemetry"
	"github.com/phanthehoang2503/small-project/order-service/internal/consumer"
	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"github.com/phanthehoang2503/small-project/order-service/internal/router"

	_ "github.com/phanthehoang2503/small-project/order-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// @title Order Service API
// @version 1.0
// @description Handles order creation, retrieval, and status updates.
// @BasePath /
// @host localhost:8083
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()
	// Init Tracer
	shutdown := telemetry.InitTracer("order-service")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}

	b := helper.ConnectRabbit()
	defer b.Close()

	// tell logger which service this is
	logger.SetService("order-service")

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	s := repo.NewOrderRepo(db)

	// Setup queue for order.paid events
	queueName := "payment_queue"
	if err := b.DeclareQueue(queueName); err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	// Bind queue to order exchange
	if err := b.BindQueue(queueName, event.ExchangeOrder, []string{event.RoutingKeyPaymentSucceeded, event.RoutingKeyPaymentFailed}); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	// Start payment consumer
	paymentConsumer := consumer.NewOrderPaidConsumer(s, b)
	if err := paymentConsumer.Start(queueName); err != nil {
		log.Fatalf("failed to start payment consumer: %v", err)
	}
	log.Println("[order-service] payment consumer started")

	// when failing
	stockQueue := "stock_failed_queue"
	if err := b.DeclareQueue(stockQueue); err != nil {
		log.Fatalf("failed to declare stock queue: %v", err)
	}
	if err := b.BindQueue(stockQueue, event.ExchangeOrder, []string{event.RoutingKeyInventoryReservationFailed}); err != nil {
		log.Fatalf("failed to bind stock queue: %v", err)
	}

	// Start stock consumer
	stockConsumer := consumer.NewStockConsumer(s, b)
	if err := stockConsumer.Start(stockQueue); err != nil {
		log.Fatalf("failed to start stock consumer: %v", err)
	}
	log.Println("[order-service] stock consumer started")

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	r := gin.Default()
	r.Use(otelgin.Middleware("order-service"))
	r.Use(middleware.CORSMiddleware())
	router.RegisterRoutes(r, s, b, jwtSecret)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8083")
}
