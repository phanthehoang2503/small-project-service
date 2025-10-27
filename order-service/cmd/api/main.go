package main

import (
	"log"

	"github.com/gin-gonic/gin"
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
	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}

	orderRepo := repo.NewOrderRepo(db)
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	router.RegisterRoutes(r, orderRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8082")
}
