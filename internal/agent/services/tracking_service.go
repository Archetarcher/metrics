package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"runtime"
)

type TrackingService struct {
	Client *resty.Client
}

func (s *TrackingService) Fetch(counterInterval int64, metrics *models.MetricsData) *models.TrackingError {
	mapMetricsValues(counterInterval, metrics)
	return nil
}

func (s *TrackingService) Send(request *models.Metrics) (*models.SendResponse, *models.TrackingError) {

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	js, err := json.Marshal(request)
	if err != nil {
		return nil, &models.TrackingError{Text: fmt.Sprintf("error marshal json: %s\n", err), Code: http.StatusInternalServerError}
	}
	_, err = zb.Write(js)
	err = zb.Close()
	if err != nil {
		return nil, &models.TrackingError{Text: fmt.Sprintf("error compression: %s\n", err), Code: http.StatusInternalServerError}
	}

	url := fmt.Sprintf("http://%s/update/", models.ServerRunAddr)

	res, err := s.Client.R().SetHeaders(models.ClientHeaders).SetBody(buf).Post(url)
	if err != nil {
		return nil, &models.TrackingError{Text: fmt.Sprintf("client: could not create request: %s\n", err.Error()), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &models.TrackingError{Text: fmt.Sprintf("client: responded with error: %s\n", err), Code: res.StatusCode()}
	}
	return &models.SendResponse{Status: http.StatusOK}, nil
}

func metricsValue(name string, mtype string, delta *int64, value *float64) models.Metrics {
	return models.Metrics{
		ID:    name,
		MType: mtype,
		Delta: delta,
		Value: value,
	}
}

func mapMetricsValues(counterInterval int64, metrics *models.MetricsData) {
	randomValue := rand.ExpFloat64()
	gauge := gatherGaugeValues()

	metrics.PollCount = metricsValue(models.PollCount, models.CounterType, &counterInterval, nil)
	metrics.RandomValue = metricsValue(models.RandomValue, models.GaugeType, nil, &randomValue)
	metrics.Alloc = metricsValue(models.Alloc, models.GaugeType, nil, &gauge.Alloc)
	metrics.BuckHashSys = metricsValue(models.BuckHashSys, models.GaugeType, nil, &gauge.BuckHashSys)
	metrics.Frees = metricsValue(models.Frees, models.GaugeType, nil, &gauge.Frees)
	metrics.GCCPUFraction = metricsValue(models.GCCPUFraction, models.GaugeType, nil, &gauge.GCCPUFraction)
	metrics.GCSys = metricsValue(models.GCSys, models.GaugeType, nil, &gauge.GCSys)
	metrics.HeapAlloc = metricsValue(models.HeapAlloc, models.GaugeType, nil, &gauge.HeapAlloc)
	metrics.HeapIdle = metricsValue(models.HeapIdle, models.GaugeType, nil, &gauge.HeapIdle)
	metrics.HeapInuse = metricsValue(models.HeapInuse, models.GaugeType, nil, &gauge.HeapInuse)
	metrics.HeapObjects = metricsValue(models.HeapObjects, models.GaugeType, nil, &gauge.HeapObjects)
	metrics.HeapReleased = metricsValue(models.HeapReleased, models.GaugeType, nil, &gauge.HeapReleased)
	metrics.HeapSys = metricsValue(models.HeapSys, models.GaugeType, nil, &gauge.HeapSys)
	metrics.LastGC = metricsValue(models.LastGC, models.GaugeType, nil, &gauge.LastGC)
	metrics.Lookups = metricsValue(models.Lookups, models.GaugeType, nil, &gauge.Lookups)
	metrics.MCacheInuse = metricsValue(models.MCacheInuse, models.GaugeType, nil, &gauge.MCacheInuse)
	metrics.MCacheSys = metricsValue(models.MCacheSys, models.GaugeType, nil, &gauge.MCacheSys)
	metrics.LastGC = metricsValue(models.LastGC, models.GaugeType, nil, &gauge.LastGC)
	metrics.Lookups = metricsValue(models.Lookups, models.GaugeType, nil, &gauge.Lookups)
	metrics.MCacheInuse = metricsValue(models.MCacheInuse, models.GaugeType, nil, &gauge.MCacheInuse)
	metrics.MCacheSys = metricsValue(models.MCacheSys, models.GaugeType, nil, &gauge.MCacheSys)
	metrics.MSpanInuse = metricsValue(models.MSpanInuse, models.GaugeType, nil, &gauge.MSpanInuse)
	metrics.MSpanSys = metricsValue(models.MSpanSys, models.GaugeType, nil, &gauge.MSpanSys)
	metrics.Mallocs = metricsValue(models.Mallocs, models.GaugeType, nil, &gauge.Mallocs)
	metrics.NextGC = metricsValue(models.NextGC, models.GaugeType, nil, &gauge.NextGC)
	metrics.NumForcedGC = metricsValue(models.NumForcedGC, models.GaugeType, nil, &gauge.NumForcedGC)
	metrics.NumGC = metricsValue(models.NumGC, models.GaugeType, nil, &gauge.NumGC)
	metrics.OtherSys = metricsValue(models.OtherSys, models.GaugeType, nil, &gauge.OtherSys)
	metrics.PauseTotalNs = metricsValue(models.PauseTotalNs, models.GaugeType, nil, &gauge.PauseTotalNs)
	metrics.StackInuse = metricsValue(models.StackInuse, models.GaugeType, nil, &gauge.StackInuse)
	metrics.StackSys = metricsValue(models.StackSys, models.GaugeType, nil, &gauge.StackSys)
	metrics.Sys = metricsValue(models.Sys, models.GaugeType, nil, &gauge.Sys)
	metrics.TotalAlloc = metricsValue(models.TotalAlloc, models.GaugeType, nil, &gauge.TotalAlloc)
}
func gatherGaugeValues() models.Gauge {
	var gauge = models.Gauge{}

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
