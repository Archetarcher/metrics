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
	FetchMemory() (*domain.MetricsData, *domain.TrackingError)
	FetchRuntime(counterInterval int64) (*domain.MetricsData, *domain.TrackingError)
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
	wg.Add(1)

	metricsData := make(chan domain.MetricsData, h.Config.RateLimit)

	for w := 1; w <= h.Config.RateLimit; w++ {
		go reportWorker(h.Send, metricsData, h.Config.ReportInterval)
	}
	go startRuntimePoll(h.FetchRuntime, &wg, h.Config.PollInterval, metricsData)
	go startMemoryPoll(h.FetchMemory, &wg, h.Config.PollInterval, metricsData)

	logger.Log.Info("Waiting for goroutines to finish...")

	wg.Wait()
	logger.Log.Info("Done!")

	return nil
}

type fetchMemory func() (*domain.MetricsData, *domain.TrackingError)
type fetchRuntime func(counterInterval int64) (*domain.MetricsData, *domain.TrackingError)
type send func(request []domain.Metrics) (*domain.SendResponse, *domain.TrackingError)

func reportWorker(send send, metricsData <-chan domain.MetricsData, interval int) {
	var reportInterval = time.Duration(interval) * time.Second
	logger.Log.Info("starting report")

	for data := range metricsData {
		vals := maps.Values(data)
		_, err := send(vals)
		if err != nil {
			logger.Log.Info(err.Text)

			retry(1, 3, vals, send)
		}
		logger.Log.Info("send data and sleep")

		time.Sleep(reportInterval)
	}
}

func startRuntimePoll(fetch fetchRuntime, wg *sync.WaitGroup, interval int, pollData chan<- domain.MetricsData) {
	defer wg.Done()
	var pollInterval = time.Duration(interval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting runtime poll")
	for {
		metrics, err := fetch(counterInterval)
		if err != nil {
			logger.Log.Info(err.Text)
		}

		pollData <- *metrics

		counterInterval++

		time.Sleep(pollInterval)
	}
}
func startMemoryPoll(fetch fetchMemory, wg *sync.WaitGroup, interval int, pollData chan<- domain.MetricsData) {
	defer wg.Done()
	var pollInterval = time.Duration(interval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting memory poll")
	for {
		metrics, err := fetch()
		if err != nil {
			logger.Log.Info(err.Text)
		}

		pollData <- *metrics

		counterInterval++

		time.Sleep(pollInterval)
	}
}

func retry(interval int, try int, vals []domain.Metrics, send send) {
	logger.Log.Info("retrying send", zap.Int("interval", interval), zap.Int("try", try))

	time.Sleep(time.Duration(interval) * time.Second)

	if try < 0 {
		return
	}

	_, err := send(vals)
	if err != nil {
		retry(interval+2, try-1, vals, send)
	}
}
