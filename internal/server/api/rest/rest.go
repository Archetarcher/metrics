package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
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

type MetricAPI struct {
	router chi.Router
}

func NewMetricAPI(storage *store.MemStorage) (*MetricAPI, error) {
	config.ParseConfig()

	r := chi.NewRouter()

	if err := logger.Initialize(domain.LogLevel); err != nil {
		return nil, err
	}
	r.Use(logger.RequestLogger)

	repo := &repositories.MetricRepository{Storage: storage}
	service := &services.MetricsService{MetricRepository: repo}
	handler := handlers.MetricsHandler{MetricsService: service}

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)
	return &MetricAPI{
		router: r,
	}, nil
}

func (a MetricAPI) Run() error {

	logger.Log.Info("Running server ", zap.String("address", domain.RunAddr))
	return http.ListenAndServe(domain.RunAddr, a.router)
}
