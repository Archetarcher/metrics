package rest

import (
	"errors"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/provider"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"go.uber.org/zap"
)

// Run starts metric tracking by rest handler
func Run(conf *config.AppConfig, s *services.MetricsService) error {
	p := provider.NewMetricsProvider(conf)

	h, err := handlers.NewMetricsHandler(conf, p, s)
	if err != nil {
		logger.Log.Info("failed to create tracking handler with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return errors.New(err.Text)
	}

	eErr := h.StartSession()
	if eErr != nil {
		logger.Log.Error("failed to start secure session", zap.String("error", eErr.Text), zap.Int("code", eErr.Code))
		return errors.New(eErr.Text)
	}

	h.TrackMetrics()

	return nil
}
