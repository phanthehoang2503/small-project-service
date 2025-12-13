package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/auth-service/internal/handler"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RegisterRoutes(r *gin.Engine, h *handler.AuthHandler, jwtSecret []byte, loginLimiter gin.HandlerFunc) {
	r.Use(otelgin.Middleware("auth-service"))
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", loginLimiter, h.Login)
	}

	// Protected API group
	api := r.Group("/api")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.GET("/profile", func(c *gin.Context) {
			uid, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": uid})
		})
	}

}
