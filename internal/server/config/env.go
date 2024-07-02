package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"os"
)

const (
	envRunAddrName  = "ADDRESS"
	envLogLevelName = "LOG_LEVEL"
)

func parseEnv() {
	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		models.RunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv(envLogLevelName); envLogLevel != "" {
		models.LogLevel = envLogLevel
	}
}
