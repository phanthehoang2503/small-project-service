package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
	"github.com/phanthehoang2503/small-project/product-service/internal/router"
)

func main() {
	r := gin.Default()
	repo := repo.NewRepo()

	router.RegisterRoutes(r, repo)
	r.Run(":8080")
}
