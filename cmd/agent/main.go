package main

import (
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func main() {
	service := &services.TrackingService{Client: resty.New()}
	handler := handlers.TrackingHandler{TrackingService: service}
	err := handler.TrackMetrics()
	if err != nil {
		logger.Log.Error("failed with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}
}
