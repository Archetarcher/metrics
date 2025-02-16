package rest

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"github.com/Archetarcher/metrics.git/internal/server/api/rest/middlewares"
	config2 "github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	service2 "github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

var conf Config

type Config struct {
	c      *config2.AppConfig
	server *httptest.Server
	err    error
	once   sync.Once
}

var counter = int64(2896127014)
var gauge = 0.31167763133187076

func (c *Config) setConfig() {
	c.once.Do(func() {
		c.c = config2.NewConfig()

		server, err := setupConfigServer()

		c.server = server
		c.err = err
	})
}
func setupConfigServer() (*httptest.Server, error) {
	ctx := context.Background()

	storage, err := store.NewStore(ctx, conf.c)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}

	repo := repositories.NewMetricsRepository(storage)
	service := service2.NewMetricsService(repo)
	handler := handlers.NewMetricsHandler(service, conf.c)
	r := chi.NewRouter()
	r.Use(middlewares.GzipMiddleware)

	r.Post("/session/", handler.StartSession)
	r.Post("/updates/", handler.UpdatesMetrics)

	srv := httptest.NewServer(r)

	return srv, nil
}

func init() {
	conf.setConfig()
}
func TestMetricsClient_Run(t *testing.T) {
	type fields struct {
		config      *config.AppConfig
		service     *services.MetricsService
		privatePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "negative test 1",
			fields: fields{
				config:      &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", ""), PublicKeyPath: "../../../../public.pem", Session: config.Session{RetryConn: 3}, RateLimit: 3},
				privatePath: "../../../private.pem",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf.c.PrivateKeyPath = tt.fields.privatePath

			c := &MetricsClient{
				config:  tt.fields.config,
				service: services.NewMetricsService(tt.fields.config),
			}
			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
