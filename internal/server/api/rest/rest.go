package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/compression"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type MetricAPI struct {
	router chi.Router
}

func NewMetricAPI(storage *store.MemStorage) (*MetricAPI, error) {
	r := chi.NewRouter()

	if err := logger.Initialize(models.LogLevel); err != nil {
		return nil, err
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
	return &MetricAPI{
		router: r,
	}, nil
}

func (a MetricAPI) Run() error {

	logger.Log.Info("Running server ", zap.String("address", models.RunAddr))
	return http.ListenAndServe(models.RunAddr, a.router)
}
