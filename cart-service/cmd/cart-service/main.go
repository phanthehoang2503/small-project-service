package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/cart-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/database"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect to database...")
	}
	cartRepo := repo.NewCartRepo(db)
	if err := db.AutoMigrate(&model.Cart{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	router.RegisterRoutes(r, cartRepo)
	r.Run(":8081")
}
