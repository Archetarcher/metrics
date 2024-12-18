package encryption

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	config2 "github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/middlewares"
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
var m = []domain.Metrics{
	{
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	},
	{
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
		ID:    "counter_value",
		MType: "counter",
		Delta: &counter,
		Value: nil,
	}, {
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
	}}

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

	srv := httptest.NewServer(r)

	return srv, nil
}

func init() {
	conf.setConfig()
}
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
			err := StartSession(tt.args.config, client, tt.args.config.Session.RetryConn)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestEncryptAsymmetric(t *testing.T) {
	require.Nil(t, conf.err, "failed to init server", conf.server, conf.err)

	type args struct {
		text []byte
		key  string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				text: []byte("teststring"),
				key:  "../../../public.pem",
			},
			wantErr: false,
		},
		{
			name: "negative test #2",
			args: args{
				text: []byte("teststring"),
				key:  "../../public.pem",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, eErr := EncryptAsymmetric(tt.args.text, tt.args.key)

			assert.Equal(t, tt.wantErr, eErr != nil)
		})
	}
}

func TestEncryptSymmetric(t *testing.T) {
	require.Nil(t, conf.err, "failed to init server", conf.server, conf.err)

	type args struct {
		text []byte
		key  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				text: []byte("teststring"),
				key:  "secretkey",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncryptSymmetric(tt.args.text, tt.args.key)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_genKey(t *testing.T) {
	require.Nil(t, conf.err, "failed to init server", conf.server, conf.err)

	tests := []struct {
		name    string
		n       int
		wantErr bool
	}{
		{
			name:    "positive test #1",
			n:       16,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := genKey(tt.n)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
