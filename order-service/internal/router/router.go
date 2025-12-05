package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/order-service/internal/handler"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RegisterRoutes(r *gin.Engine, s *repo.OrderRepo, b *broker.Broker, jwtSecret []byte) {
	r.Use(otelgin.Middleware("order-service"))
	api := r.Group("/orders")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.POST("", handler.CreateOrder(s, b)) // List all orders for a specific user (?user_id=1)
		api.GET("", handler.ListOrders(s))
		api.GET("/search", handler.SearchOrders(s)) // Search order by ID (?id=1)
		api.GET("/:id", handler.GetOrder(s))
		api.PUT("/:id/status", handler.UpdateOrderStatus(s)) // Update order status (/orders/:id/status)
	}
}
