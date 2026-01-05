package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/auth-service/internal/handler"
	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	"github.com/phanthehoang2503/small-project/auth-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"github.com/phanthehoang2503/small-project/internal/telemetry"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	_ "github.com/phanthehoang2503/small-project/auth-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// @title Auth Service API
// @version 1.0
// @description Authentication service (register / login / token)
// @host localhost:8084
// @BasePath /
func main() {
	godotenv.Load()
	// Init Tracer
	shutdown := telemetry.InitTracer("auth-service")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	db := DbConnect()

	// connect to rabbit
	b := helper.ConnectRabbit()
	defer b.Close()

	// tell logger which service this is
	logger.SetService("auth-service")

	userRepo := repo.NewUserRepo(db)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	authHandler := handler.NewAuthHandler(userRepo, jwtSecret, 72)

	// Connect Redis for Rate Limiting
	redisAddr := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	r := gin.Default()
	r.Use(otelgin.Middleware("auth-service"))
	r.Use(middleware.CORSMiddleware())

	// Apply rate limit specifically to login
	rateLimitStr := os.Getenv("LOGIN_RATE_LIMIT")
	limit := 10 // default
	if val, err := strconv.Atoi(rateLimitStr); err == nil && val > 0 {
		limit = val
	}
	loginLimiter := middleware.RateLimitMiddleware(rdb, limit, time.Minute)

	router.RegisterRoutes(r, authHandler, jwtSecret, loginLimiter)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8084")
}

func DbConnect() *gorm.DB {
	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect to database...")
	}
	return db
}
