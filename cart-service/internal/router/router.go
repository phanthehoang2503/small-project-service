package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/handler"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
)

func RegisterRoutes(r *gin.Engine, s *repo.CartRepo) {
	api := r.Group("/cart")
	{
		api.POST("", handler.AddToCart(s))
	}
}
