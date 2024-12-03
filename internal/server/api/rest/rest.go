package rest

import (
	"context"
	"errors"
	"github.com/Archetarcher/metrics.git/internal/server/encryption"
	"github.com/Archetarcher/metrics.git/internal/server/middlewares"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

// MetricsAPI is an api struct, keeps router.
type MetricsAPI struct {
	router chi.Router
}

// NewMetricsAPI registers routes, middlewares.
func NewMetricsAPI(handler *handlers.MetricsHandler, config *config.AppConfig) (*MetricsAPI, *domain.MetricsError) {

	r := chi.NewRouter()
	r.Use(func(handler http.Handler) http.Handler {
		return encryption.RequestDecryptMiddleware(handler, config)
	})
	r.Use(middlewares.GzipMiddleware)

	r.Use(logger.RequestLoggerMiddleware)

	r.Use(func(handler http.Handler) http.Handler {
		return middlewares.RequestHashesMiddleware(handler, config)
	})

	r.Mount("/debug", middleware.Profiler())

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)
	r.Get("/value/{type}/{name}", handler.GetMetrics)
	r.Get("/", handler.GetMetricsPage)

	r.Post("/update/", handler.UpdateMetricsJSON)
	r.Post("/updates/", handler.UpdatesMetrics)
	r.Post("/value/", handler.GetMetricsJSON)
	r.Post("/session/", handler.StartSession)

	r.Get("/ping", handler.GetPing)
	return &MetricsAPI{
		router: r,
	}, nil
}

// RunRestServer starts serving application.
func RunRestServer(config *config.AppConfig, repo *repositories.MetricRepository) error {
	service := services.NewMetricsService(repo)
	handler := handlers.NewMetricsHandler(service, config)

	api, err := NewMetricsAPI(handler, config)

	if err != nil {
		logger.Log.Error("failed to init api with error, finishing app", zap.String("error", err.Text), zap.Int("code", err.Code))
		return err.Err
	}

	logger.Log.Info("Running rest server ", zap.String("address", config.RunAddr))

	server := &http.Server{Addr: config.RunAddr, Handler: api.router}
	configShutdown(server)

	if lErr := server.ListenAndServe(); !errors.Is(lErr, http.ErrServerClosed) {
		logger.Log.Error("failed to serve rest server", zap.String("error", err.Text), zap.Int("code", err.Code))
		return lErr
	}
	return nil
}

func configShutdown(srv *http.Server) {
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		<-sigint
		logger.Log.Info("got interruption signal")
		time.Sleep(time.Duration(10) * time.Second)

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Info("HTTP server Shutdown: ", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

}
