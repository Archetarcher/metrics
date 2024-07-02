package handlers

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"sync"
	"time"
)

type TrackingHandler struct {
	TrackingService
}

type TrackingService interface {
	Fetch(counterInterval int64) ([]models.Metrics, *models.TrackingError)
	Send(request *models.Metrics) (*models.SendResponse, *models.TrackingError)
}

func (h *TrackingHandler) TrackMetrics() error {
	config.ParseConfig()

	if err := logger.Initialize(models.LogLevel); err != nil {
		return err
	}

	fmt.Println("start tracking")

	var wg sync.WaitGroup
	metrics := make(chan models.Metrics)
	wg.Add(2)

	go startPoll(h.Fetch, metrics, &wg)
	go startReport(h.Send, metrics, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
	return nil
}

type fetch func(counterInterval int64) ([]models.Metrics, *models.TrackingError)
type send func(request *models.Metrics) (*models.SendResponse, *models.TrackingError)

func startPoll(fetch fetch, metrics chan<- models.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	var pollInterval = time.Duration(models.PollInterval) * time.Second
	counterInterval := int64(1)
	fmt.Println("starting poll")
	for {
		response, err := fetch(counterInterval)
		if err != nil {
			logger.Log.Error(err.Text)
		}

		for _, m := range response {
			fmt.Println("write to chan")
			fmt.Println(m)
			metrics <- m
		}
		counterInterval++

		time.Sleep(pollInterval)
	}
}

func startReport(send send, metrics <-chan models.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("starting report")
	var reportInterval = time.Duration(models.ReportInterval) * time.Second
	for metric := range metrics {
		fmt.Println("reading from chan")
		fmt.Println(metric)

		response, err := send(&metric)
		if err != nil {
			logger.Log.Error(err.Text)
		}
		fmt.Println(response)

		time.Sleep(reportInterval)
	}
}
