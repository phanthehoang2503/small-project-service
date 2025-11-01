package router

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/auth-service/internal/handler"
	"github.com/phanthehoang2503/small-project/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, h *handler.AuthHandler, jwtSecret []byte) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}

	// Protected API group — reuse shared JWT middleware from pkg/auth
	api := r.Group("/api")
	api.Use(middleware.JWTMiddleware(jwtSecret))
	{
		api.GET("/profile", func(c *gin.Context) {
			uid, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": uid})
		})
	}

}
