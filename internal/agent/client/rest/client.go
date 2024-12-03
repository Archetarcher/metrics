package rest

import (
	"errors"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"go.uber.org/zap"
)

// Run starts metric tracking by rest handler
func Run(conf *config.AppConfig) error {
	h, err := handlers.NewMetricsHandler(conf)
	if err != nil {
		logger.Log.Info("failed to create tracking handler with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return errors.New(err.Text)
	}

	eErr := h.StartSession()
	if eErr != nil {
		logger.Log.Error("failed to start secure session", zap.String("error", eErr.Text), zap.Int("code", eErr.Code))
		return errors.New(eErr.Text)
	}

	hErr := h.TrackMetrics()
	if hErr != nil {
		logger.Log.Info("failed with error", zap.String("error", hErr.Text), zap.Int("code", hErr.Code))
		return errors.New(hErr.Text)
	}
	return nil
}
