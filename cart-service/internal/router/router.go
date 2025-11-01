package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/handler"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, s *repo.CartRepo, jwtSecret []byte) {
	api := r.Group("/cart")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.POST("", handler.AddToCart(s))
		api.GET("", handler.GetCart(s))
		api.PUT("/:id", handler.UpdateQuantity(s))
		api.DELETE("/:id", handler.RemoveItem(s))
		api.DELETE("", handler.ClearCart(s))
	}
}
