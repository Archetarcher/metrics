package repositories

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
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

var conf Config

type Config struct {
	c    *config.AppConfig
	repo *MetricRepository
	err  error
	once sync.Once
}

func (c *Config) setConfig() {
	c.once.Do(func() {
		c.c = config.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

		repo, err := setup()

		c.repo = repo
		c.err = err
	})
}
func init() {
	conf.setConfig()
}

func setup() (*MetricRepository, error) {
	ctx := context.Background()

	storage, err := store.NewStore(ctx, conf.c.Store)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
		return nil, err.Err
	}

	repo := NewMetricsRepository(storage)

	return repo, nil
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

func TestMetricRepository_Set(t *testing.T) {
	require.NoError(t, conf.err, "failed to init repo", conf.repo, conf.err)

	type args struct {
		request *domain.Metrics
	}

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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.repo.Set(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestMetricRepository_Get(t *testing.T) {
	require.NoError(t, conf.err, "failed to init repo", conf.repo, conf.err)

	type args struct {
		request *domain.Metrics
	}

	tests := []struct {
		want    *domain.Metrics
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
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.repo.Get(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
func TestMetricRepository_GetAll(t *testing.T) {
	require.NoError(t, conf.err, "failed to init repo", conf.repo, conf.err)

	type args struct {
	}

	tests := []struct {
		want    *domain.Metrics
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.repo.GetAll(ctx)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricRepository_GetAllIn(t *testing.T) {
	require.NoError(t, conf.err, "failed to init repo", conf.repo, conf.err)

	tests := []struct {
		want    *domain.Metrics
		args    []string
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
			args:    []string{values[0].ID + values[0].MType},
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.repo.GetAllIn(ctx, tt.args)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricRepository_SetAll(t *testing.T) {
	require.NoError(t, conf.err, "failed to init repo", conf.repo, conf.err)

	type args struct {
		request []domain.Metrics
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: []domain.Metrics{values[1], values[2]}},
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.repo.SetAll(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestMetricRepository_CheckConnection(t *testing.T) {
	require.NoError(t, conf.err, "failed to init service", conf.repo, conf.err)

	t.Run("test check connection", func(t *testing.T) {
		ctx := context.Background()
		err := conf.repo.CheckConnection(ctx)
		assert.Nil(t, err, "CheckConnection(%v)")
	})
}

func TestNewMetricsRepository(t *testing.T) {

	t.Run("positive test", func(t *testing.T) {
		m := NewMetricsRepository(&pgx.Store{})
		assert.NotNil(t, m)
	})
}
