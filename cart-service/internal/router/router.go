package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/handler"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RegisterRoutes(r *gin.Engine, cartRepo repo.CartRepository, productRepo *repo.ProductRepo, jwtSecret []byte) {
	r.Use(otelgin.Middleware("cart-service"))
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
