package services

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/agent/types"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	gaugeType     = "gauge"
	counterType   = "counter"
	pollCount     = "PollCount"
	randomValue   = "RandomValue"
	alloc         = "Alloc"
	buckHashSys   = "BuckHashSys"
	frees         = "Frees"
	gCCPUFraction = "GCCPUFraction"
	gCSys         = "GCSys"
	heapAlloc     = "HeapAlloc"
	heapIdle      = "HeapIdle"
	heapInuse     = "HeapInuse"
	heapObjects   = "HeapObjects"
	heapReleased  = "HeapReleased"
	heapSys       = "HeapSys"
	lastGC        = "LastGC"
	lookups       = "Lookups"
	mCacheInuse   = "MCacheInuse"
	mCacheSys     = "MCacheSys"
	mSpanInuse    = "MSpanInuse"
	mSpanSys      = "MSpanSys"
	mallocs       = "Mallocs"
	nextGC        = "NextGC"
	numForcedGC   = "NumForcedGC"
	numGC         = "NumGC"
	otherSys      = "OtherSys"
	pauseTotalNs  = "PauseTotalNs"
	stackInuse    = "StackInuse"
	stackSys      = "StackSys"
	sys           = "Sys"
	totalAlloc    = "TotalAlloc"
)

// MetricsService is a service for gathering and sending metrics to server.
type MetricsService struct {
	config *config.AppConfig
}

func NewMetricsService(appConfig *config.AppConfig) *MetricsService {
	return &MetricsService{config: appConfig}
}

// TrackMetrics fetches runtime gauge and counter metrics.
func (s *MetricsService) TrackMetrics(ctx context.Context, update types.UpdateMetrics, group *sync.WaitGroup) {

	metricsData := make(chan domain.MetricsData, s.config.RateLimit)

	for w := 1; w <= s.config.RateLimit; w++ {
		go reportWorker(update, metricsData, time.Duration(s.config.ReportInterval)*time.Second, w)
	}
	go startRuntimePoll(group, s.config.PollInterval, metricsData, ctx)
	go startMemoryPoll(group, s.config.PollInterval, metricsData, ctx)
}

func reportWorker(send types.UpdateMetrics, metricsData <-chan domain.MetricsData, reportInterval time.Duration, index int) {

	logger.Log.Info("starting report")

	for d := range metricsData {
		logger.Log.Info("Worker started metric ", zap.Int("worker_id", index))

		vals := maps.Values(d)
		_, err := send(vals)
		if err != nil {
			logger.Log.Info(err.Text)

			retry(1, 3, vals, send)
		}
		logger.Log.Info("Worker finished processing ", zap.Int("worker_id", index))

		time.Sleep(reportInterval)
	}
}

func startRuntimePoll(wg *sync.WaitGroup, interval int, pollData chan domain.MetricsData, ctx context.Context) {
	defer wg.Done()
	pollInterval := time.Duration(interval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting runtime poll")
	for {
		select {
		case <-ctx.Done():
			if !isClosed(pollData) {
				close(pollData)
			}
			return
		default:
			metrics, err := fetchRuntime(counterInterval)
			if err != nil {
				logger.Log.Info(err.Text)
			}

			pollData <- *metrics

			counterInterval++

			time.Sleep(pollInterval)
		}

	}
}
func startMemoryPoll(wg *sync.WaitGroup, interval int, pollData chan domain.MetricsData, ctx context.Context) {
	defer wg.Done()
	pollInterval := time.Duration(interval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting memory poll")
	for {
		select {
		case <-ctx.Done():
			if !isClosed(pollData) {
				close(pollData)
			}
			return
		default:
			metrics, err := fetchMemory()
			if err != nil {
				logger.Log.Info(err.Text)
			}

			pollData <- *metrics

			counterInterval++

			time.Sleep(pollInterval)
		}
	}
}

func retry(interval int, try int, vals []domain.Metrics, send types.UpdateMetrics) {
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

func fetchMemory() (*domain.MetricsData, *domain.MetricsError) {
	metrics := mapGaugeMetrics(gatherMemoryValues)
	return &metrics, nil
}
func fetchRuntime(counterInterval int64) (*domain.MetricsData, *domain.MetricsError) {
	metrics := mapGaugeMetrics(gatherRuntimeValues)
	metrics[pollCount] = metricsValue(pollCount, counterType, &counterInterval, nil)

	return &metrics, nil
}

func metricsValue(name string, mtype string, delta *int64, value *float64) domain.Metrics {
	return domain.Metrics{
		ID:    name,
		MType: mtype,
		Delta: delta,
		Value: value,
	}
}

func mapGaugeMetrics(value gatherGaugeValue) domain.MetricsData {
	rv := rand.ExpFloat64()
	gauge := value()
	metrics := make(map[string]domain.Metrics)

	metrics[randomValue] = metricsValue(randomValue, gaugeType, nil, &rv)
	metrics[alloc] = metricsValue(alloc, gaugeType, nil, &gauge.Alloc)
	metrics[buckHashSys] = metricsValue(buckHashSys, gaugeType, nil, &gauge.BuckHashSys)
	metrics[frees] = metricsValue(frees, gaugeType, nil, &gauge.Frees)
	metrics[gCCPUFraction] = metricsValue(gCCPUFraction, gaugeType, nil, &gauge.GCCPUFraction)
	metrics[gCSys] = metricsValue(gCSys, gaugeType, nil, &gauge.GCSys)
	metrics[heapAlloc] = metricsValue(heapAlloc, gaugeType, nil, &gauge.HeapAlloc)
	metrics[heapIdle] = metricsValue(heapIdle, gaugeType, nil, &gauge.HeapIdle)
	metrics[heapInuse] = metricsValue(heapInuse, gaugeType, nil, &gauge.HeapInuse)
	metrics[heapObjects] = metricsValue(heapObjects, gaugeType, nil, &gauge.HeapObjects)
	metrics[heapReleased] = metricsValue(heapReleased, gaugeType, nil, &gauge.HeapReleased)
	metrics[heapSys] = metricsValue(heapSys, gaugeType, nil, &gauge.HeapSys)
	metrics[lastGC] = metricsValue(lastGC, gaugeType, nil, &gauge.LastGC)
	metrics[lookups] = metricsValue(lookups, gaugeType, nil, &gauge.Lookups)
	metrics[mCacheInuse] = metricsValue(mCacheInuse, gaugeType, nil, &gauge.MCacheInuse)
	metrics[mCacheSys] = metricsValue(mCacheSys, gaugeType, nil, &gauge.MCacheSys)
	metrics[lastGC] = metricsValue(lastGC, gaugeType, nil, &gauge.LastGC)
	metrics[lookups] = metricsValue(lookups, gaugeType, nil, &gauge.Lookups)
	metrics[mCacheInuse] = metricsValue(mCacheInuse, gaugeType, nil, &gauge.MCacheInuse)
	metrics[mCacheSys] = metricsValue(mCacheSys, gaugeType, nil, &gauge.MCacheSys)
	metrics[mSpanInuse] = metricsValue(mSpanInuse, gaugeType, nil, &gauge.MSpanInuse)
	metrics[mSpanSys] = metricsValue(mSpanSys, gaugeType, nil, &gauge.MSpanSys)
	metrics[mallocs] = metricsValue(mallocs, gaugeType, nil, &gauge.Mallocs)
	metrics[nextGC] = metricsValue(nextGC, gaugeType, nil, &gauge.NextGC)
	metrics[numForcedGC] = metricsValue(numForcedGC, gaugeType, nil, &gauge.NumForcedGC)
	metrics[numGC] = metricsValue(numGC, gaugeType, nil, &gauge.NumGC)
	metrics[otherSys] = metricsValue(otherSys, gaugeType, nil, &gauge.OtherSys)
	metrics[pauseTotalNs] = metricsValue(pauseTotalNs, gaugeType, nil, &gauge.PauseTotalNs)
	metrics[stackInuse] = metricsValue(stackInuse, gaugeType, nil, &gauge.StackInuse)
	metrics[stackSys] = metricsValue(stackSys, gaugeType, nil, &gauge.StackSys)
	metrics[sys] = metricsValue(sys, gaugeType, nil, &gauge.Sys)
	metrics[totalAlloc] = metricsValue(totalAlloc, gaugeType, nil, &gauge.TotalAlloc)
	return metrics
}

func gatherRuntimeValues() domain.Gauge {
	var gauge domain.Gauge
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	gauge.Alloc = float64(rtm.Alloc)
	gauge.BuckHashSys = float64(rtm.BuckHashSys)
	gauge.Frees = float64(rtm.Frees)
	gauge.GCCPUFraction = rtm.GCCPUFraction
	gauge.GCSys = float64(rtm.GCSys)
	gauge.HeapAlloc = float64(rtm.HeapAlloc)
	gauge.HeapIdle = float64(rtm.HeapIdle)
	gauge.HeapInuse = float64(rtm.HeapInuse)
	gauge.HeapObjects = float64(rtm.HeapObjects)
	gauge.HeapReleased = float64(rtm.HeapReleased)
	gauge.HeapSys = float64(rtm.HeapSys)
	gauge.LastGC = float64(rtm.LastGC)
	gauge.Lookups = float64(rtm.Lookups)
	gauge.MCacheInuse = float64(rtm.MCacheInuse)
	gauge.MCacheSys = float64(rtm.MCacheSys)
	gauge.MSpanInuse = float64(rtm.MSpanInuse)
	gauge.MSpanSys = float64(rtm.MSpanSys)
	gauge.Mallocs = float64(rtm.Mallocs)
	gauge.NextGC = float64(rtm.NextGC)
	gauge.NumForcedGC = float64(rtm.NumForcedGC)
	gauge.NumGC = float64(rtm.NumGC)
	gauge.OtherSys = float64(rtm.OtherSys)
	gauge.PauseTotalNs = float64(rtm.PauseTotalNs)
	gauge.StackInuse = float64(rtm.StackInuse)
	gauge.StackSys = float64(rtm.StackSys)
	gauge.Sys = float64(rtm.Sys)
	gauge.TotalAlloc = float64(rtm.TotalAlloc)
	gauge.RandomValue = rand.ExpFloat64()
	return gauge
}
func gatherMemoryValues() domain.Gauge {
	var gauge domain.Gauge
	vm, _ := mem.VirtualMemory()

	gauge.TotalMemory = float64(vm.Total)
	gauge.FreeMemory = float64(vm.Free)

	cp, _ := cpu.PercentWithContext(context.TODO(), 0, true)
	gauge.CPUutilization1 = float64(len(cp))

	return gauge
}

type gatherGaugeValue func() domain.Gauge

func isClosed(ch <-chan domain.MetricsData) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
