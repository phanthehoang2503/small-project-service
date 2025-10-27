package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/cart-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/database"

	_ "github.com/phanthehoang2503/small-project/cart-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cart Service API
// @version 1.0
// @description Manage shopping cart items
// @host localhost:8081
// @BasePath /
func main() {
	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect to database...")
	}
	cartRepo := repo.NewCartRepo(db)
	if err := db.AutoMigrate(&model.Cart{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	godotenv.Load()

	r := gin.Default()
	router.RegisterRoutes(r, cartRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")
}
