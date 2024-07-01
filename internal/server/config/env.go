package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"os"
)

const (
	envRunAddrName = "ADDRESS"
)

func parseEnv() {
	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		domain.RunAddr = envRunAddr
	}
}
