package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type API struct {
	router chi.Router
}

func NewAPI(storage *store.MemStorage) API {
	r := chi.NewRouter()

	//mux := http.NewServeMux()
	repo := &repositories.MetricRepository{Storage: storage}
	service := &services.MetricsService{MetricRepositoryInterface: repo}
	handler := handlers.MetricsHandler{MetricsServiceInterface: service}

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)
	return API{
		router: r,
	}
}

func (a API) Run() error {
	return http.ListenAndServe(`:8080`, a.router)
}
