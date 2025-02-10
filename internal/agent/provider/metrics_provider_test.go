package provider

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/Archetarcher/metrics.git/internal/server/api/rest/middlewares"
	config2 "github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	service := services.NewMetricsService(repo)
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

var confClient *config.AppConfig

func TestStartSession(t *testing.T) {
	require.Nil(t, conf.err, "failed to init server", conf.server, conf.err)

	type args struct {
		config      *config.AppConfig
		privatePath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				config:      &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", ""), PublicKeyPath: "../../../public.pem", Session: config.Session{RetryConn: 3}},
				privatePath: "../../../private.pem",
			},
			wantErr: false,
		},
		{
			name: "negative test #2",
			args: args{
				config:      &config.AppConfig{ServerRunAddr: conf.server.URL, PublicKeyPath: "../../../public.pem", Session: config.Session{RetryConn: 3}},
				privatePath: "../../../private.pem",
			},
			wantErr: true,
		},
		{
			name: "negative test #3",
			args: args{
				config:      &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", ""), PublicKeyPath: "../../public.pem", Session: config.Session{RetryConn: 3}},
				privatePath: "../../../private.pem",
			},
			wantErr: true,
		},
		{
			name: "negative test #4",
			args: args{
				config:      &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", ""), PublicKeyPath: "../../../public.pem", Session: config.Session{RetryConn: 2}},
				privatePath: "../../private.pem",
			},
			wantErr: true,
		},
	}

	client := resty.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf.c.PrivateKeyPath = tt.args.privatePath
			prvdr := NewMetricsProvider(tt.args.config, client)

			err := prvdr.StartSession(tt.args.config.Session.RetryConn)
			if err == nil {
				confClient = tt.args.config
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

var (
	values = [8]domain.Metrics{
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

//func TestUpdate(t *testing.T) {
//
//	require.Nil(t, conf.err, "failed to init server", conf.server, conf.err)
//
//	type args struct {
//		request []domain.Metrics
//	}
//
//	tests := []struct {
//		name    string
//		args    args
//		code    int
//		config  *config.AppConfig
//		wantErr bool
//	}{
//		{
//			name: "positive test #1",
//			args: args{[]domain.Metrics{values[0], values[1]}},
//			code: http.StatusOK,
//			config: &config.AppConfig{ServerRunAddr: strings.ReplaceAll(conf.server.URL, "http://", ""),
//				Session: confClient.Session, PublicKeyPath: confClient.PublicKeyPath},
//			wantErr: false,
//		},
//		{
//			name: "negative test #2",
//			args: args{[]domain.Metrics{
//				{
//					ID:    values[0].ID,
//					MType: "gauged",
//					Value: nil,
//				},
//			},
//			},
//			config: &config.AppConfig{ServerRunAddr: confClient.ServerRunAddr,
//				Session: confClient.Session, PublicKeyPath: confClient.PublicKeyPath},
//			code:    http.StatusBadRequest,
//			wantErr: true,
//		},
//		{
//			name: "negative test #3",
//			args: args{[]domain.Metrics{values[0]}},
//			code: http.StatusInternalServerError,
//			config: &config.AppConfig{ServerRunAddr: conf.server.URL,
//				Session: confClient.Session, PublicKeyPath: confClient.PublicKeyPath},
//			wantErr: true,
//		},
//	}
//	client := resty.New()
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			prvdr := NewMetricsProvider(confClient, client)
//
//			_, pErr := prvdr.Update(tt.args.request)
//
//			assert.Equal(t, tt.wantErr, pErr != nil)
//		})
//	}
//}
