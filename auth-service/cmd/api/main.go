package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/auth-service/internal/handler"
	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	"github.com/phanthehoang2503/small-project/auth-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/database"
	"github.com/phanthehoang2503/small-project/internal/helper"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"gorm.io/gorm"

	_ "github.com/phanthehoang2503/small-project/auth-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description Authentication service (register / login / token)
// @host localhost:8084
// @BasePath /
func main() {
	_ = godotenv.Load()

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

	r := gin.Default()
	router.RegisterRoutes(r, authHandler, jwtSecret)

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
