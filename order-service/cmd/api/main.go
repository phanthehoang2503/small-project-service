package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"github.com/phanthehoang2503/small-project/order-service/internal/router"

	_ "github.com/phanthehoang2503/small-project/order-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Order Service API
// @version 1.0
// @description API documentation for the order microservice
// @BasePath /
func main() {
	godotenv.Load()
	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	orderRepo := repo.NewOrderRepo(db)

	r := gin.Default()
	router.RegisterRoutes(r, orderRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8082")
}
