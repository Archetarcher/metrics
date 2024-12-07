package services

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/compression"
	config2 "github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
)

var (
	counter = int64(2896127014)
	gauge   = 0.31167763133187076
	values  = [8]domain.Metrics{
		{
			ID:    "counter_value",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_2",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_2",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_3",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_3",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_4",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_4",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
	}
)

func TestTrackingService_FetchMemory(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.FetchMemory()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTrackingService_FetchRuntime(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		counter int64
	}{
		{
			name:    "positive test #1",
			wantErr: false,
			counter: int64(1),
		},
	}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.FetchRuntime(tt.counter)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

var conf Config

type Config struct {
	c      *config2.AppConfig
	server *httptest.Server
	err    error
	once   sync.Once
}

func (c *Config) setConfig() {
	c.once.Do(func() {
		c.c = config2.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

		server, err := setupConfigServer()

		c.server = server
		c.err = err
	})
}
func setupConfigServer() (*httptest.Server, error) {
	ctx := context.Background()

	storage, err := store.NewStore(ctx, conf.c.Store)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}

	repo := repositories.NewMetricsRepository(storage)
	service := services.NewMetricsService(repo)
	handler := handlers.NewMetricsHandler(service, conf.c)
	r := chi.NewRouter()
	r.Use(compression.GzipMiddleware)

	r.Post("/updates/", handler.UpdatesMetrics)

	srv := httptest.NewServer(r)

	return srv, nil
}

func init() {
	conf.setConfig()
}

func TestTrackingService_Send(t *testing.T) {

	type args struct {
		request []domain.Metrics
	}
	tests := []struct {
		name    string
		args    args
		code    int
		config  *config.AppConfig
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{[]domain.Metrics{values[0]}},
			code:    http.StatusOK,
			config:  &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", "")},
			wantErr: false,
		},
		{
			name: "positive test #3",
			args: args{[]domain.Metrics{
				{
					ID:    values[0].ID,
					MType: "gauged",
					Value: nil,
				},
			},
			},
			config:  &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", "")},
			code:    http.StatusBadRequest,
			wantErr: true,
		},
		{
			name:    "positive test #3",
			args:    args{[]domain.Metrics{values[0]}},
			code:    http.StatusInternalServerError,
			config:  &config.AppConfig{ServerRunAddr: conf.server.URL},
			wantErr: true,
		},
	}

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TrackingService{Config: tt.config, Client: client}

			_, err := service.Send(tt.args.request)
			assert.Equal(t, tt.wantErr, err != nil, err)

			if err != nil {
				assert.Equal(t, tt.code, err.Code, err)

			}

		})
	}
}
