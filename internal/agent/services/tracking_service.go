package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
)

type TrackingService struct {
	Client *resty.Client
}

func (s *TrackingService) Fetch(counterInterval int64) ([]models.Metrics, *models.TrackingError) {

	var metrics = make([]models.Metrics, 0)

	var gauge = mapMetricValues()

	metrics = append(metrics, models.Metrics{
		ID:    models.CounterMetric,
		MType: models.CounterType,
		Delta: &counterInterval,
	})
	values := reflect.ValueOf(gauge)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		name := types.Field(i).Name
		field := values.FieldByName(name).Float()

		metrics = append(metrics, models.Metrics{
			ID:    name,
			MType: models.GaugeType,
			Value: &field,
		})
	}

	return metrics, nil
}

func (s *TrackingService) Send(request *models.Metrics) (*models.SendResponse, *models.TrackingError) {

	url := fmt.Sprintf("http://%s/update/", models.ServerRunAddr)
	logger.LogData("send update request ", request)

	res, err := s.Client.R().SetHeader("Content-Type", "Content-Type: application/json").SetBody(request).Post(url)
	if err != nil {
		return nil, &models.TrackingError{Text: fmt.Sprintf("client: could not create request: %s\n", err), Code: http.StatusInternalServerError}
	}

	if res.StatusCode() != http.StatusOK {
		return nil, &models.TrackingError{Text: fmt.Sprintf("client: responded with error: %s\n", err), Code: res.StatusCode()}
	}
	return &models.SendResponse{Status: http.StatusOK}, nil
}

func mapMetricValues() models.Gauge {
	var rtm runtime.MemStats
	var gauge models.Gauge

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
