package rest

import (
	"errors"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/provider"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type MetricsClient struct {
	config  *config.AppConfig
	service *services.MetricsService
}

func NewMetricsClient(config *config.AppConfig, service *services.MetricsService) *MetricsClient {
	return &MetricsClient{config: config, service: service}
}

// Run starts metric tracking by rest handler
func (c *MetricsClient) Run() error {
	client := resty.New()

	p := provider.NewMetricsProvider(c.config, client)
	h, err := handlers.NewMetricsHandler(c.config, p, c.service)
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
