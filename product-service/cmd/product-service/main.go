package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/product-service/internal/handler"
	"github.com/phanthehoang2503/product-service/internal/store"
)

func main() {
	r := gin.Default()
	store := store.NewStore()

	handler.RegisterRoutes(r, store)
	r.Run(":8080")
}
