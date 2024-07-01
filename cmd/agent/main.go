package main

import (
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"github.com/go-resty/resty/v2"
)

func main() {
	service := &services.TrackingService{Client: resty.New()}
	handler := handlers.TrackingHandler{TrackingService: service}
	handler.TrackMetrics()
}
