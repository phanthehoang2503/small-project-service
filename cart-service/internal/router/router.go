package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/handler"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, cartRepo *repo.CartRepo, productRepo *repo.ProductRepo, jwtSecret []byte) {
	api := r.Group("/cart")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.POST("", handler.AddToCart(cartRepo, productRepo))
		api.GET("", handler.GetCart(cartRepo))
		api.PUT("/:id", handler.UpdateQuantity(cartRepo))
		api.DELETE("/:id", handler.RemoveItem(cartRepo))
		api.DELETE("", handler.ClearCart(cartRepo))
	}
}
