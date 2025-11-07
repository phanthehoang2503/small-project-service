package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/logger-service/internal/handler"
	"github.com/phanthehoang2503/logger-service/internal/logger"
	"github.com/rs/zerolog"
)

func main() {
	// init zerolog
	zlog := logger.InitLogger("logger-service", "./logs/app.log", zerolog.DebugLevel)

	r := gin.Default()
	h := handler.NewHandler(zlog)
	r.POST("/ingest", h.ReceiveLog)
	r.Run(":8085")
}
