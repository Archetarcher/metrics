package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"os"
)

const (
	envRunAddrName  = "ADDRESS"
	envLogLevelName = "LOG_LEVEL"
)

func parseEnv() {
	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		domain.RunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv(envLogLevelName); envLogLevel != "" {
		domain.LogLevel = envLogLevel
	}
}
