package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/product-service/internal/handler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
)

func RegisterRoutes(r *gin.Engine, s *repo.Database, cache *repo.CacheRepository) {
	r.Use(otelgin.Middleware("product-service"))
	api := r.Group("/products")
	{
		api.GET("", handler.ListProducts(s))
		api.GET("/:id", handler.GetProducts(s, cache))
		api.POST("", handler.CreateProducts(s))
		api.PUT("/:id", handler.UpdateProducts(s, cache))
		api.DELETE("/:id", handler.DeleteProducts(s, cache))
	}
}
