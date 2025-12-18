package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/internal/telemetry"
	_ "github.com/phanthehoang2503/small-project/product-service/docs"
	"github.com/phanthehoang2503/small-project/product-service/internal/consumer"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"github.com/phanthehoang2503/small-project/product-service/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// @title Product Service API
// @version 1.0
// @description Manage product items
// @host localhost:8081
// @BasePath /
func main() {
	godotenv.Load()
	// Init Tracer
	shutdown := telemetry.InitTracer("product-service")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	//DB
	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	productRepo := repo.NewRepo(db)

	// RabbitMQ + lging
	b := helper.ConnectRabbit()
	defer b.Close()

	logger.SetService("product-service")

	// Redis Cache
	cacheRepo := repo.NewCacheRepository("redis:6379")

	// Setup Queue & Binding
	queueName := "product_order_events"
	if err := b.DeclareQueue(queueName); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	// Bind to Order Exchange
	bindingKeys := []string{
		event.RoutingKeyOrderCreated,
		event.RoutingKeyPaymentFailed,
	}
	if err := b.BindQueue(queueName, event.ExchangeOrder, bindingKeys); err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// Start Consumer
	orderConsumer := consumer.NewOrderConsumer(productRepo, cacheRepo, b)
	if err := orderConsumer.Start(queueName); err != nil {
		log.Printf("Failed to start consumer: %v", err)
	}

	r := gin.Default()
	r.Use(otelgin.Middleware("product-service"))
	r.Use(middleware.CORSMiddleware())
	router.RegisterRoutes(r, productRepo, cacheRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")
}
