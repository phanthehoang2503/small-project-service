package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/product-service/internal/handler"

	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
)

func RegisterRoutes(r *gin.Engine, s *repo.Database) {
	api := r.Group("/products")
	{
		api.GET("", handler.ListProducts(s))
		api.GET("/:id", handler.GetProducts(s))
		api.POST("", handler.CreateProducts(s))
		api.PUT("/:id", handler.UpdateProducts(s))
		api.DELETE("/:id", handler.DeleteProducts(s))
	}
}
