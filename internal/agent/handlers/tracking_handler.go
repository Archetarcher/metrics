package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"net/http"
	"sync"
	"time"
)

type TrackingHandler struct {
	TrackingService
	Config *config.AppConfig
}

type TrackingService interface {
	Fetch(counterInterval int64, metrics *domain.MetricsData) *domain.TrackingError
	Send(request []domain.Metrics) (*domain.SendResponse, *domain.TrackingError)
}

func (h *TrackingHandler) TrackMetrics() *domain.TrackingError {

	if err := logger.Initialize(h.Config.LogLevel); err != nil {
		return &domain.TrackingError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
	logger.Log.Info("start tracking")

	var wg sync.WaitGroup
	metrics := domain.MetricsData{}
	wg.Add(2)

	go startPoll(h.Fetch, &metrics, &wg, h.Config.PollInterval)
	go startReport(h.Send, &metrics, &wg, h.Config.ReportInterval)

	logger.Log.Info("Waiting for goroutines to finish...")

	wg.Wait()
	logger.Log.Info("Done!")
	return nil
}

type fetch func(counterInterval int64, metrics *domain.MetricsData) *domain.TrackingError
type send func(request []domain.Metrics) (*domain.SendResponse, *domain.TrackingError)

func startPoll(fetch fetch, metrics *domain.MetricsData, wg *sync.WaitGroup, interval int) {
	defer wg.Done()
	var pollInterval = time.Duration(interval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting poll")
	for {
		err := fetch(counterInterval, metrics)
		if err != nil {
			logger.Log.Info(err.Text)
		}

		counterInterval++

		time.Sleep(pollInterval)
	}
}

func startReport(send send, metrics *domain.MetricsData, wg *sync.WaitGroup, interval int) {
	defer wg.Done()

	logger.Log.Info("starting report")

	var reportInterval = time.Duration(interval) * time.Second
	for {

		vals := maps.Values(*metrics)
		_, err := send(vals)
		if err != nil {
			logger.Log.Info(err.Text)

			retry(1, 3, vals, send)
		}

		time.Sleep(reportInterval)

	}

}

func retry(interval int, try int, vals []domain.Metrics, send send) {
	logger.Log.Info("retrying send", zap.Int("interval", interval), zap.Int("try", try))

	time.Sleep(time.Duration(interval) * time.Second)

	if try > 0 {
		return
	}

	_, err := send(vals)
	if err != nil {
		retry(interval+2, try-1, vals, send)
	}
}
