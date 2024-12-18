package rest

import (
	"context"
	"errors"
	"github.com/Archetarcher/metrics.git/internal/server/api/rest/middlewares"
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
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

// MetricsServer is an api struct, keeps router.
type MetricsServer struct {
	router  chi.Router
	handler *handlers.MetricsHandler
	config  *config.AppConfig
}

func NewMetricsServer(service *services.MetricsService, config *config.AppConfig) *MetricsServer {
	handler := handlers.NewMetricsHandler(service, config)

	return &MetricsServer{
		router:  chi.NewRouter(),
		config:  config,
		handler: handler,
	}
}

// Run starts serving application.
func (a *MetricsServer) Run() error {
	logger.Log.Info("Running rest server ", zap.String("address", a.config.RunAddr))
	a.mountMiddleware()
	a.mountRoutes()

	server := &http.Server{Addr: a.config.RunAddr, Handler: a.router}
	configShutdown(server)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error("failed to serve rest server")
		return err
	}
	return nil
}
func (a *MetricsServer) mountRoutes() {
	a.router.Mount("/debug", middleware.Profiler())

	a.router.Post("/update/{type}/{name}/{value}", a.handler.UpdateMetrics)
	a.router.Get("/value/{type}/{name}", a.handler.GetMetrics)
	a.router.Get("/", a.handler.GetMetricsPage)

	a.router.Post("/update/", a.handler.UpdateMetricsJSON)
	a.router.Post("/updates/", a.handler.UpdatesMetrics)
	a.router.Post("/value/", a.handler.GetMetricsJSON)
	a.router.Post("/session/", a.handler.StartSession)

	a.router.Get("/ping", a.handler.GetPing)
}

func (a *MetricsServer) mountMiddleware() {
	a.router.Use(func(handler http.Handler) http.Handler {
		return middlewares.RequestTrustedSubnet(handler, a.config)
	})
	a.router.Use(func(handler http.Handler) http.Handler {
		return middlewares.RequestDecryptMiddleware(handler, a.config)
	})
	a.router.Use(middlewares.GzipMiddleware)
	a.router.Use(middlewares.RequestLoggerMiddleware)
	a.router.Use(func(handler http.Handler) http.Handler {
		return middlewares.RequestHashesMiddleware(handler, a.config)
	})

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
