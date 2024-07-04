package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/compression"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type MetricsAPI struct {
	router chi.Router
}

func NewMetricsAPI(storage *store.MemStorage) (*MetricsAPI, *domain.MetricsError) {
	r := chi.NewRouter()

	if err := logger.Initialize(domain.LogLevel); err != nil {
		return nil, &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
	r.Use(compression.GzipMiddleware)
	r.Use(logger.RequestLoggerMiddleware)

	repo := &repositories.MetricRepository{Storage: storage}
	service := &services.MetricsService{MetricRepository: repo}
	handler := handlers.MetricsHandler{MetricsService: service}

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)

	r.Post("/update/", handler.UpdateMetricsJSON)
	r.Post("/value/", handler.GetMetricsJSON)
	return &MetricsAPI{
		router: r,
	}, nil
}

func (a MetricsAPI) Run() *domain.MetricsError {

	logger.Log.Info("Running server ", zap.String("address", domain.RunAddr))
	err := http.ListenAndServe(domain.RunAddr, a.router)
	if err != nil {
		return &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}
