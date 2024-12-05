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
func NewMetricsHandler(conf *config.AppConfig, provider MetricsProvider, service MetricsService) (*MetricsHandler, *domain.MetricsError) {
	return &MetricsHandler{service: service, config: conf, provider: provider}, nil
}

// StartSession is a handler that generates session
func (h *MetricsHandler) StartSession() *domain.MetricsError {
	return h.provider.StartSession(h.config.Session.RetryConn)
}

// TrackMetrics starts metrics tracking.
// Runs worker pool reportWorker with domain.MetricsData chanel, each worker sends data to server.
// Runs two goroutines for runtime and memory metrics, each goroutine pulls data to domain.MetricsData chanel.
func (h *MetricsHandler) TrackMetrics() *domain.MetricsError {

	logger.Log.Info("start tracking")

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	h.service.TrackMetrics(ctx, h.provider.Update, &wg)
	logger.Log.Info("Waiting for goroutines to finish...")

	configShutdown(cancelFunc, &wg)
	logger.Log.Info("Done!")

	//
	//metricsData := make(chan domain.MetricsData, h.Config.RateLimit)
	//
	//for w := 1; w <= h.Config.RateLimit; w++ {
	//	go reportWorker(h.provider.Update, metricsData, time.Duration(h.Config.ReportInterval)*time.Second, w)
	//}
	//go startRuntimePoll(h.service.FetchRuntime, &wg, h.Config.PollInterval, metricsData, ctx)
	//go startMemoryPoll(h.service.FetchMemory, &wg, h.Config.PollInterval, metricsData, ctx)

	return nil
}

func configShutdown(cancelFunc context.CancelFunc, group *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigint

	logger.Log.Info("Shutdown signal received")

	cancelFunc()
	group.Wait()
}

//type fetchMemory func() (*domain.MetricsData, *domain.MetricsError)
//type fetchRuntime func(counterInterval int64) (*domain.MetricsData, *domain.MetricsError)

//func reportWorker(send updateMetrics, metricsData <-chan domain.MetricsData, reportInterval time.Duration, index int) {
//
//	logger.Log.Info("starting report")
//
//	for d := range metricsData {
//		logger.Log.Info("Worker started metric ", zap.Int("worker_id", index))
//
//		vals := maps.Values(d)
//		_, err := send(vals)
//		if err != nil {
//			logger.Log.Info(err.Text)
//
//			retry(1, 3, vals, send)
//		}
//		logger.Log.Info("Worker finished processing ", zap.Int("worker_id", index))
//
//		time.Sleep(reportInterval)
//	}
//}
//
//func startRuntimePoll(fetch fetchRuntime, wg *sync.WaitGroup, interval int, pollData chan<- domain.MetricsData, ctx context.Context) {
//	defer wg.Done()
//	pollInterval := time.Duration(interval) * time.Second
//	counterInterval := int64(1)
//	logger.Log.Info("starting runtime poll")
//	for {
//		select {
//		case <-ctx.Done():
//			close(pollData)
//			return
//		default:
//			metrics, err := fetch(counterInterval)
//			if err != nil {
//				logger.Log.Info(err.Text)
//			}
//
//			pollData <- *metrics
//
//			counterInterval++
//
//			time.Sleep(pollInterval)
//		}
//
//	}
//}
//func startMemoryPoll(fetch fetchMemory, wg *sync.WaitGroup, interval int, pollData chan<- domain.MetricsData, ctx context.Context) {
//	defer wg.Done()
//	pollInterval := time.Duration(interval) * time.Second
//	counterInterval := int64(1)
//	logger.Log.Info("starting memory poll")
//	for {
//		select {
//		case <-ctx.Done():
//			close(pollData)
//			return
//		default:
//			metrics, err := fetch()
//			if err != nil {
//				logger.Log.Info(err.Text)
//			}
//
//			pollData <- *metrics
//
//			counterInterval++
//
//			time.Sleep(pollInterval)
//		}
//	}
//}

//func retry(interval int, try int, vals []domain.Metrics, send updateMetrics) {
//	logger.Log.Info("retrying send", zap.Int("interval", interval), zap.Int("try", try))
//
//	time.Sleep(time.Duration(interval) * time.Second)
//
//	if try < 0 {
//		return
//	}
//
//	_, err := send(vals)
//	if err != nil {
//		retry(interval+2, try-1, vals, send)
//	}
//}
