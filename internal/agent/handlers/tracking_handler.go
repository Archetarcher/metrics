package handlers

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"reflect"
	"sync"
	"time"
)

type TrackingHandler struct {
	TrackingService
}

type TrackingService interface {
	Fetch(counterInterval int64, metrics *models.MetricsData) *models.TrackingError
	Send(request *models.Metrics) (*models.SendResponse, *models.TrackingError)
}

func (h *TrackingHandler) TrackMetrics() error {
	config.ParseConfig()

	if err := logger.Initialize(models.LogLevel); err != nil {
		return err
	}
	logger.Log.Info("start tracking")

	var wg sync.WaitGroup
	metrics := models.MetricsData{}
	wg.Add(2)

	go startPoll(h.Fetch, &metrics, &wg)
	go startReport(h.Send, &metrics, &wg)

	logger.Log.Info("Waiting for goroutines to finish...")

	wg.Wait()
	logger.Log.Info("Done!")
	return nil
}

type fetch func(counterInterval int64, metrics *models.MetricsData) *models.TrackingError
type send func(request *models.Metrics) (*models.SendResponse, *models.TrackingError)

func startPoll(fetch fetch, metrics *models.MetricsData, wg *sync.WaitGroup) {
	defer wg.Done()
	var pollInterval = time.Duration(models.PollInterval) * time.Second
	counterInterval := int64(1)
	logger.Log.Info("starting poll")
	for {
		err := fetch(counterInterval, metrics)
		if err != nil {
			logger.Log.Error(err.Text)
		}

		counterInterval++

		time.Sleep(pollInterval)
	}
}

func startReport(send send, metrics *models.MetricsData, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Log.Info("starting report")

	var reportInterval = time.Duration(models.ReportInterval) * time.Second
	for {

		values := reflect.ValueOf(metrics).Elem()
		for i := 0; i < values.NumField(); i++ {
			field := values.Field(i)

			request := field.Interface().(models.Metrics)
			fmt.Println("request")
			fmt.Println(request)
			_, err := send(&request)
			if err != nil {
				logger.Log.Error(err.Text)
			}
		}
		time.Sleep(reportInterval)

	}

}
