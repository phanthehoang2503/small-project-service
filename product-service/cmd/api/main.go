package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/database"
	_ "github.com/phanthehoang2503/small-project/product-service/docs"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"github.com/phanthehoang2503/small-project/product-service/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Product Service API
// @version 1.0
// @description Manage product items
// @host localhost:8080
// @BasePath /
func main() {
	db, err := database.ConnectDB() //connect to db
	if err != nil {
		log.Fatal("failed to connect to database...")
	}
	productRepo := repo.NewRepo(db) //initial repository
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	router.RegisterRoutes(r, productRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8081")
}
