package services

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sync"
	"testing"
)

var conf Config

type Config struct {
	once    sync.Once
	c       *config.AppConfig
	service *MetricsService
	err     error
}

func (c *Config) setConfig() {
	c.once.Do(func() {
		c.c = config.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

		service, err := setup()

		c.service = service
		c.err = err
	})
}
func init() {
	conf.setConfig()
}

func setup() (*MetricsService, error) {
	ctx := context.Background()
	storage, err := store.NewStore(conf.c.Store, ctx)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}

	repo := repositories.NewMetricsRepository(storage)
	service := NewMetricsService(repo)
	return service, nil
}

var counter = int64(2896127014)
var gauge = 0.31167763133187076
var values = [8]domain.Metrics{
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

func TestMetricsService_Update(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	type args struct {
		request *domain.Metrics
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: &values[0]},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.service.Update(tt.args.request, ctx)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricsService_Updates(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	type args struct {
		request []domain.Metrics
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: []domain.Metrics{values[0], values[1], values[2]}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.service.Updates(tt.args.request, ctx)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricsService_GetValue(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	type args struct {
		request *domain.Metrics
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: &values[0]},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.service.GetValue(tt.args.request, ctx)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricsService_GetAllValues(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.service.GetAllValues(ctx)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
