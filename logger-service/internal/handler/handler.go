package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/logger-service/internal/model"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) ReceiveLog(c *gin.Context) {
	var le model.LogEntry
	if err := c.ShouldBindJSON(&le); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad payload: " + err.Error()})
		return
	}

	if le.Timestamp.IsZero() {
		le.Timestamp = time.Now().UTC()
	}
	if le.Service == "" {
		le.Service = "unknown"
	}
	if le.Level == "" {
		le.Level = "info"
	}

	event := h.Logger.With().
		Str("trace_id", le.TraceID).
		Str("service", le.Service).
		Str("level", le.Level).
		Time("timestamp", le.Timestamp).
		Logger()

	// attach fields map (if present) as nested fields
	if le.Fields != nil {
		for k, v := range le.Fields {
			event = event.With().Interface(k, v).Logger()
		}
	}

	switch le.Level {
	case "debug":
		event.Debug().Msg(le.Message)
	case "warn", "warning":
		event.Warn().Msg(le.Message)
	case "error":
		event.Error().Msg(le.Message)
	case "fatal":
		event.Fatal().Msg(le.Message)
	default:
		event.Info().Msg(le.Message)
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
