package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"runtime"
)

var clientHeaders = map[string]string{
	"Content-Type":     "Content-Type: application/json",
	"Content-Encoding": "gzip",
}

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

type TrackingService struct {
	Client *resty.Client
}

func (s *TrackingService) Fetch(counterInterval int64, metrics *domain.MetricsData) *domain.TrackingError {
	var m domain.MetricsData
	mapMetricsValues(counterInterval, &m)
	metrics = &m
	return nil
}

func (s *TrackingService) Send(request *domain.Metrics) (*domain.SendResponse, *domain.TrackingError) {

	url := fmt.Sprintf("http://%s/update/", domain.ServerRunAddr)

	body, cErr := compress(request)
	if cErr != nil {
		return nil, cErr
	}

	res, err := s.Client.R().SetHeaders(clientHeaders).SetBody(body).Post(url)
	if err != nil {
		return nil, &domain.TrackingError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &domain.TrackingError{Text: fmt.Sprintf("client: responded with error: %s\n", err), Code: res.StatusCode()}
	}
	return &domain.SendResponse{Status: http.StatusOK}, nil
}

func metricsValue(name string, mtype string, delta *int64, value *float64) domain.Metrics {
	return domain.Metrics{
		ID:    name,
		MType: mtype,
		Delta: delta,
		Value: value,
	}
}

func mapMetricsValues(counterInterval int64, metrics *domain.MetricsData) {
	rv := rand.ExpFloat64()
	gauge := gatherGaugeValues()

	metrics.PollCount = metricsValue(pollCount, counterType, &counterInterval, nil)
	metrics.RandomValue = metricsValue(randomValue, gaugeType, nil, &rv)
	metrics.Alloc = metricsValue(alloc, gaugeType, nil, &gauge.Alloc)
	metrics.BuckHashSys = metricsValue(buckHashSys, gaugeType, nil, &gauge.BuckHashSys)
	metrics.Frees = metricsValue(frees, gaugeType, nil, &gauge.Frees)
	metrics.GCCPUFraction = metricsValue(gCCPUFraction, gaugeType, nil, &gauge.GCCPUFraction)
	metrics.GCSys = metricsValue(gCSys, gaugeType, nil, &gauge.GCSys)
	metrics.HeapAlloc = metricsValue(heapAlloc, gaugeType, nil, &gauge.HeapAlloc)
	metrics.HeapIdle = metricsValue(heapIdle, gaugeType, nil, &gauge.HeapIdle)
	metrics.HeapInuse = metricsValue(heapInuse, gaugeType, nil, &gauge.HeapInuse)
	metrics.HeapObjects = metricsValue(heapObjects, gaugeType, nil, &gauge.HeapObjects)
	metrics.HeapReleased = metricsValue(heapReleased, gaugeType, nil, &gauge.HeapReleased)
	metrics.HeapSys = metricsValue(heapSys, gaugeType, nil, &gauge.HeapSys)
	metrics.LastGC = metricsValue(lastGC, gaugeType, nil, &gauge.LastGC)
	metrics.Lookups = metricsValue(lookups, gaugeType, nil, &gauge.Lookups)
	metrics.MCacheInuse = metricsValue(mCacheInuse, gaugeType, nil, &gauge.MCacheInuse)
	metrics.MCacheSys = metricsValue(mCacheSys, gaugeType, nil, &gauge.MCacheSys)
	metrics.LastGC = metricsValue(lastGC, gaugeType, nil, &gauge.LastGC)
	metrics.Lookups = metricsValue(lookups, gaugeType, nil, &gauge.Lookups)
	metrics.MCacheInuse = metricsValue(mCacheInuse, gaugeType, nil, &gauge.MCacheInuse)
	metrics.MCacheSys = metricsValue(mCacheSys, gaugeType, nil, &gauge.MCacheSys)
	metrics.MSpanInuse = metricsValue(mSpanInuse, gaugeType, nil, &gauge.MSpanInuse)
	metrics.MSpanSys = metricsValue(mSpanSys, gaugeType, nil, &gauge.MSpanSys)
	metrics.Mallocs = metricsValue(mallocs, gaugeType, nil, &gauge.Mallocs)
	metrics.NextGC = metricsValue(nextGC, gaugeType, nil, &gauge.NextGC)
	metrics.NumForcedGC = metricsValue(numForcedGC, gaugeType, nil, &gauge.NumForcedGC)
	metrics.NumGC = metricsValue(numGC, gaugeType, nil, &gauge.NumGC)
	metrics.OtherSys = metricsValue(otherSys, gaugeType, nil, &gauge.OtherSys)
	metrics.PauseTotalNs = metricsValue(pauseTotalNs, gaugeType, nil, &gauge.PauseTotalNs)
	metrics.StackInuse = metricsValue(stackInuse, gaugeType, nil, &gauge.StackInuse)
	metrics.StackSys = metricsValue(stackSys, gaugeType, nil, &gauge.StackSys)
	metrics.Sys = metricsValue(sys, gaugeType, nil, &gauge.Sys)
	metrics.TotalAlloc = metricsValue(totalAlloc, gaugeType, nil, &gauge.TotalAlloc)
}
func gatherGaugeValues() domain.Gauge {
	var gauge = domain.Gauge{}

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

func compress(request *domain.Metrics) (*bytes.Buffer, *domain.TrackingError) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	js, err := json.Marshal(request)
	if err != nil {
		return nil, &domain.TrackingError{Text: fmt.Sprintf("error marshal json: %s\n", err), Code: http.StatusInternalServerError}
	}
	_, err = zb.Write(js)

	if err != nil {
		return nil, &domain.TrackingError{Text: fmt.Sprintf("error compression: %s\n", err), Code: http.StatusInternalServerError}
	}
	err = zb.Close()
	if err != nil {
		return nil, &domain.TrackingError{Text: fmt.Sprintf("error compression: %s\n", err), Code: http.StatusInternalServerError}
	}
	return buf, nil
}
