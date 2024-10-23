package main

import (
	"context"
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

	_ "net/http/pprof"
)

func main() {
	c := config.NewConfig(store.Config{Memory: &memory.Config{}, Pgx: &pgx.Config{}})
	c.ParseConfig()

	if err := logger.Initialize(c.LogLevel); err != nil {
		log.Fatal("failed to init logger")
	}
	ctx := context.Background()

	storage, err := store.NewStore(c.Store, ctx)

	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Error(err.Err))

		ns, e := store.Retry(err, 1, 3, c.Store, ctx)

		if e != nil {
			logger.Log.Error("failed to retry init storage with error, finishing app", zap.String("error", e.Text), zap.Error(e.Err))
			return
		}
		storage = ns
	}

	repo := repositories.NewMetricsRepository(storage)
	service := services.NewMetricsService(repo)
	handler := handlers.NewMetricsHandler(service, c)

	api, err := rest.NewMetricsAPI(handler, c)

	if err != nil {
		logger.Log.Error("failed to init api with error, finishing app", zap.String("error", err.Text), zap.Int("code", err.Code))
		return
	}

	if err := api.Run(c); err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}
