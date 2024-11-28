package services

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

var conf Config

type Config struct {
	c       *config.AppConfig
	service *MetricsService
	err     error
	once    sync.Once
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
	storage, err := store.NewStore(ctx, conf.c.Store)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}

	repo := repositories.NewMetricsRepository(storage)
	service := NewMetricsService(repo)
	return service, nil
}

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

func TestMetricsService_Update(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	type args struct {
		request *domain.Metrics
	}
	ctx := context.Background()

	tests := []struct {
		args    args
		name    string
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
		args    args
		name    string
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

func TestMetricsService_CheckConnection(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.service, conf.err)

	t.Run("test check connection", func(t *testing.T) {
		ctx := context.Background()
		err := conf.service.CheckConnection(ctx)
		assert.Nil(t, err, "CheckConnection(%v)")
	})
}

func TestNewMetricsService(t *testing.T) {
	t.Run("positive test", func(t *testing.T) {
		s := NewMetricsService(&repositories.MetricRepository{})
		assert.NotNil(t, s)
	})
}

func BenchmarkMetricsService_Updates(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()
	metrics := []domain.Metrics{values[0], values[1], values[2], values[3]}
	for i := 0; i < b.N; i++ {
		conf.service.Updates(metrics, ctx)
	}
}
