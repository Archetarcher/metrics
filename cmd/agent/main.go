package main

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildData()

	conf := config.NewConfig()
	conf.ParseConfig()
	client := resty.New()

	eErr := encryption.StartSession(conf, client)
	if eErr != nil {
		logger.Log.Error("failed to start secure session", zap.String("error", eErr.Text), zap.Int("code", eErr.Code))
		return
	}
	service := &services.TrackingService{Client: client, Config: conf}
	handler := handlers.TrackingHandler{TrackingService: service, Config: conf}
	err := handler.TrackMetrics()
	if err != nil {
		logger.Log.Info("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}

func printBuildData() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
