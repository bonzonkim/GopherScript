package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

func NewLogger(env string) *Logger {
	var config zap.Config
	if env == "dev" {
		config = zap.NewDevelopmentConfig()
	}
	if env == "prod" {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger : %v", err)
	}

	return &Logger{
		Logger: logger,
	}
}
