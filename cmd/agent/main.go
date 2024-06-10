package main

import (
	"github.com/Archetarcher/metrics.git/internal/agent/handlers"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
)

func main() {
	run()
}

func run() {
	service := &services.TrackingService{}
	handler := handlers.TrackingHandler{TrackingServiceInterface: service}
	handler.StartTracking()
}
