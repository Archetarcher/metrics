package handlers

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"sync"
	"time"
)

type TrackingHandler struct {
	TrackingService
}

type TrackingService interface {
	Fetch(counterInterval int) ([]domain.MetricData, *domain.TrackingError)
	Send(request *domain.MetricData) (*domain.SendResponse, *domain.TrackingError)
}

func (h *TrackingHandler) TrackMetrics() {
	config.ParseConfig()

	fmt.Println("start tracking")

	var wg sync.WaitGroup
	metrics := make(chan domain.MetricData)
	wg.Add(2)

	go startPoll(h.Fetch, metrics, &wg)
	go startReport(h.Send, metrics, &wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
}

type fetch func(counterInterval int) ([]domain.MetricData, *domain.TrackingError)
type send func(request *domain.MetricData) (*domain.SendResponse, *domain.TrackingError)

func startPoll(fetch fetch, metrics chan<- domain.MetricData, wg *sync.WaitGroup) {
	defer wg.Done()
	var pollInterval = time.Duration(domain.PollInterval) * time.Second
	counterInterval := 1
	fmt.Println("starting poll")
	for {
		response, err := fetch(counterInterval)
		if err != nil {
			fmt.Println(err)
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

func startReport(send send, metrics <-chan domain.MetricData, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("starting report")

	var reportInterval = time.Duration(domain.ReportInterval) * time.Second

	for metric := range metrics {
		fmt.Println("reading from chan")
		fmt.Println(metric)

		response, err := send(&metric)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response)

		time.Sleep(reportInterval)
	}
}
