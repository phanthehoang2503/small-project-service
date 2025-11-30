package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	_ "github.com/phanthehoang2503/small-project/product-service/docs"
	"github.com/phanthehoang2503/small-project/product-service/internal/consumer"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"github.com/phanthehoang2503/small-project/product-service/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Product Service API
// @version 1.0
// @description Manage product items
// @host localhost:8081
// @BasePath /
func main() {
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

	// Start Consumer
	orderConsumer := consumer.NewOrderConsumer(productRepo, b)
	if err := orderConsumer.Start("product_order_events"); err != nil {
		log.Printf("Failed to start consumer: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.RegisterRoutes(r, productRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")
}
