package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"github.com/phanthehoang2503/small-project/order-service/internal/router"

	_ "github.com/phanthehoang2503/small-project/order-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}

	b := helper.ConnectRabbit()
	defer b.Close()

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	orderRepo := repo.NewOrderRepo(db)
	secret := []byte(os.Getenv("JWT_SECRET"))
	r := gin.Default()
	router.RegisterRoutes(r, orderRepo, b, secret)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8083")
}
