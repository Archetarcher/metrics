package handlers

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/types"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// MetricsHandler is a handler for tracking metrics, has service and configuration.
type MetricsHandler struct {
	service  MetricsService
	provider MetricsProvider
	config   *config.AppConfig
}

// MetricsService is an interface for tracking metrics, sends and fetch memory and runtime metrics.
type MetricsService interface {
	TrackMetrics(ctx context.Context, update types.UpdateMetrics, group *sync.WaitGroup)
}

// MetricsProvider is an interface for sending metrics to server
type MetricsProvider interface {
	Update(request []domain.Metrics) (*domain.SendResponse, *domain.MetricsError)
	StartSession(retryCount int) *domain.MetricsError
}

// NewMetricsHandler creates and sets up tracking handler
func NewMetricsHandler(conf *config.AppConfig, provider MetricsProvider, service MetricsService) *MetricsHandler {
	return &MetricsHandler{service: service, config: conf, provider: provider}
}

// StartSession is a handler that generates session
func (h *MetricsHandler) StartSession() *domain.MetricsError {
	return h.provider.StartSession(h.config.Session.RetryConn)
}

// TrackMetrics starts metrics tracking.
// Runs worker pool reportWorker with domain.MetricsData chanel, each worker sends data to server.
// Runs two goroutines for runtime and memory metrics, each goroutine pulls data to domain.MetricsData chanel.
func (h *MetricsHandler) TrackMetrics() {

	logger.Log.Info("start tracking")

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	h.service.TrackMetrics(ctx, h.provider.Update, &wg)
	logger.Log.Info("Waiting for goroutines to finish...")

	configShutdown(cancelFunc, &wg)
	logger.Log.Info("Done!")

}

func configShutdown(cancelFunc context.CancelFunc, group *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigint

	logger.Log.Info("Shutdown signal received")

	cancelFunc()
	group.Wait()
}
