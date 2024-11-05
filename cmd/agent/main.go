package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

var (
	buildVersion = ""
	buildDate    = ""
	buildCommit  = ""
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	c := config.NewConfig()
	c.ParseConfig()
	service := &services.TrackingService{Client: resty.New(), Config: c}
	handler := handlers.TrackingHandler{TrackingService: service, Config: c}
	err := handler.TrackMetrics()
	if err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}
