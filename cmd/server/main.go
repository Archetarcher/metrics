package main

import (
	"context"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"go.uber.org/zap"
	"log"
	_ "net/http/pprof"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildData()

	c := config.NewConfig()
	c.ParseConfig()

	if err := logger.Initialize(c.LogLevel); err != nil {
		log.Fatal("failed to init logger")
	}
	ctx := context.Background()

	storage, err := store.NewStore(ctx, c)

	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Error(err.Err))

		ns, e := store.Retry(ctx, err, 1, 3, c)

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
	aErr := api.Run(c)
	if aErr != nil {
		logger.Log.Error("failed to start server with error, finishing app", zap.Error(aErr))
		return
	}
}

func printBuildData() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
