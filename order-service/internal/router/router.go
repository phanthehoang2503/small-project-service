package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/order-service/internal/handler"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
)

func RegisterRoutes(r *gin.Engine, s *repo.OrderRepo, jwtSecret []byte) {
	api := r.Group("/orders")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.POST("", handler.CreateOrder(s)) // List all orders for a specific user (?user_id=1)
		api.GET("", handler.ListOrders(s))
		api.GET("/:id", handler.GetOrder(s))
		api.PUT("/:id/status", handler.UpdateOrderStatus(s)) // Update order status (/orders/:id/status)
	}
}
