package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type MetricAPI struct {
	router chi.Router
}

func NewAPI(storage *store.MemStorage) MetricAPI {
	config.ParseConfig()

	r := chi.NewRouter()

	repo := &repositories.MetricRepository{Storage: storage}
	service := &services.MetricsService{MetricRepository: repo}
	handler := handlers.MetricsHandler{MetricsService: service}

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)
	return MetricAPI{
		router: r,
	}
}

func (a MetricAPI) Run() error {
	return http.ListenAndServe(domain.RunAddr, a.router)
}
