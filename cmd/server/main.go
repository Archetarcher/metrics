package main

import (
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"go.uber.org/zap"
	"log"
)

func main() {
	c := config.NewConfig(store.Config{Memory: &memory.Config{}, Pgx: &pgx.Config{}})
	c.ParseConfig()

	if err := logger.Initialize(c.LogLevel); err != nil {
		log.Fatal("failed to init logger")
	}

	storage, err := store.NewStore(c.Store)

	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	repo := repositories.NewMetricsRepository(storage)
	service := services.NewMetricsService(repo)
	handler := handlers.NewMetricsHandler(service, c)

	api, err := rest.NewMetricsAPI(handler)

	if err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	if err := api.Run(c); err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}
