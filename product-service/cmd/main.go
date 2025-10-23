package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"github.com/phanthehoang2503/small-project/product-service/internal/router"
)

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
	r.Run(":8080")
}
