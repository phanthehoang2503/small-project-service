package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(serviceName string, logFile string, level zerolog.Level) zerolog.Logger {
	lumber := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     14,   // days
		Compress:   true, // gzip
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// combine writers: console + file
	multi := io.MultiWriter(consoleWriter, lumber)

	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(multi).With().Timestamp().Str("service", serviceName).Logger().Level(level)

	return logger
}
