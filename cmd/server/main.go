package main

import (
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"go.uber.org/zap"
)

func main() {
	c := config.NewConfig()
	c.ParseConfig()

	storage := store.NewStorage(c)

	api, err := rest.NewMetricsAPI(storage, c)

	if err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	if err := api.Run(); err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}
