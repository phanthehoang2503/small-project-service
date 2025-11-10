package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phanthehoang2503/small-project/auth-service/internal/handler"
	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	"github.com/phanthehoang2503/small-project/auth-service/internal/router"
	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/database"
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

	rabbitURL := os.Getenv("RABBITMQ_URL")
	var b *broker.Broker
	var err error

	for i := 0; i < 10; i++ {
		b, err = broker.InitRabbit(rabbitURL, "logs_exchange")
		if err == nil {
			log.Println("connected to RabbitMQ")
			break
		}
		log.Printf("attempt %d: failed to connect to RabbitMQ (%v)\n", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("could not connect to RabbitMQ after multiple attempts:", err)
	}
	defer b.Close()
	log.Println("RabbitMQ ready in auth-service")

	userRepo := repo.NewUserRepo(db)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	authHandler := handler.NewAuthHandler(userRepo, jwtSecret, 72)

	r := gin.Default()
	router.RegisterRoutes(r, authHandler, jwtSecret)

	logger.SetConfig("", "auth-service")

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
