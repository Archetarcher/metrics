package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/compression"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/encoding"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type MetricsAPI struct {
	router chi.Router
}

func NewMetricsAPI(handler *handlers.MetricsHandler, config *config.AppConfig) (*MetricsAPI, *domain.MetricsError) {
	r := chi.NewRouter()

	r.Use(compression.GzipMiddleware)
	r.Use(logger.RequestLoggerMiddleware)
	r.Use(func(handler http.Handler) http.Handler {
		return encoding.RequestHashesMiddleware(handler, config)
	})

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)

	r.Post("/update/", handler.UpdateMetricsJSON)
	r.Post("/updates/", handler.UpdatesMetrics)
	r.Post("/value/", handler.GetMetricsJSON)

	r.Get("/ping", handler.GetPing)
	return &MetricsAPI{
		router: r,
	}, nil
}

func (a MetricsAPI) Run(config *config.AppConfig) *domain.MetricsError {

	logger.Log.Info("Running server ", zap.String("address", config.RunAddr))
	err := http.ListenAndServe(config.RunAddr, a.router)
	if err != nil {
		return &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}
