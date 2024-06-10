package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
)

type API struct {
	server *http.ServeMux
}

func NewAPI(storage *store.MemStorage) API {

	mux := http.NewServeMux()
	repo := &repositories.MetricRepository{Storage: storage}
	service := &services.MetricsService{MetricRepositoryInterface: repo}
	handler := handlers.MetricsHandler{MetricsServiceInterface: service}

	mux.HandleFunc(`/update/{type}/{name}/{value}`, handler.UpdateMetrics)

	return API{
		server: mux,
	}
}

func (a API) Run() error {
	return http.ListenAndServe(`:8080`, a.server)
}
