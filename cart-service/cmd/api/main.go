package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/cart-service/internal/consumer"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/cart-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/event"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "github.com/phanthehoang2503/small-project/cart-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cart Service API
// @version 1.0
// @description Manage shopping cart items
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()

	// Init Tracer
	shutdown := telemetry.InitTracer("cart-service")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect to database...")
	}

	b := helper.ConnectRabbit()
	defer b.Close()

	// Migrations
	if err := db.AutoMigrate(&model.Cart{}); err != nil {
		log.Fatalf("Migration failed (cart): %v", err)
	}
	if err := db.AutoMigrate(&model.ProductSnapshot{}); err != nil {
		log.Fatalf("Migration failed (product_snapshot): %v", err)
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	//repos
	cr := repo.NewCartRepo(db)
	pr := repo.NewProductRepo(db)

	// Setup Product Queue
	prodQueue := "cart_products_queue"
	if err := b.DeclareQueue(prodQueue); err != nil {
		log.Fatalf("failed to declare product queue: %v", err)
	}
	if err := b.BindQueue(prodQueue, event.ExchangeProduct, []string{"product.*"}); err != nil {
		log.Fatalf("failed to bind product queue: %v", err)
	}

	pc := consumer.NewProductConsumer(pr)
	if err := pc.Start(event.ExchangeProduct, prodQueue, "product.*"); err != nil {
		log.Fatalf("failed to start product consumer: %v", err)
	}

	// Setup Order Queue
	orderQueue := "cart_orders_queue"
	if err := b.DeclareQueue(orderQueue); err != nil {
		log.Fatalf("failed to declare order queue: %v", err)
	}
	if err := b.BindQueue(orderQueue, event.ExchangeOrder, []string{event.RoutingKeyOrderRequested}); err != nil {
		log.Fatalf("failed to bind order queue: %v", err)
	}

	oc := consumer.NewOrderConsumer(cr, b)
	if err := oc.Start(event.ExchangeOrder, orderQueue, event.RoutingKeyOrderRequested); err != nil {
		log.Fatalf("failed to start order consumer: %v", err)
	}

	r := gin.Default()
	r.Use(otelgin.Middleware("cart-service"))
	r.Use(middleware.CORSMiddleware())
	router.RegisterRoutes(r, cr, pr, jwtSecret)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8082")
}
