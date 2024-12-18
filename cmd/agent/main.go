package main

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildData()

	h, err := handlers.NewTrackingHandler()
	if err != nil {
		logger.Log.Info("failed to create tracking handler with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return
	}

	hErr := h.TrackMetrics()
	if hErr != nil {
		logger.Log.Info("failed with error", zap.String("error", hErr.Text), zap.Int("code", hErr.Code))
		return
	}
}

func printBuildData() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
