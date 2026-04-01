package helpers

import (
	"tiny-goclean/config"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(cfg *config.Config) *zap.Logger {
	if cfg.Server.Environment == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	// TODO: define later
	defConfig := &config.Config{}
	return InitLogger(defConfig)
}
