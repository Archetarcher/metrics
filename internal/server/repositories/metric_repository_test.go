package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var c = config.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

func setup() (*domain.Metrics, *MetricRepository) {
	i := int64(1)
	req := &domain.Metrics{MType: "counter", ID: "countervalue", Delta: &i}

	storage, err := store.NewStore(c.Store)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	repo := NewMetricsRepository(storage)

	return req, repo
}
func TestMetricRepository_Get(t *testing.T) {

	c.ParseConfig()
	type args struct {
		request *domain.Metrics
	}

	req, repo := setup()
	tests := []struct {
		name    string
		args    args
		want    *domain.Metrics
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: req},
			wantErr: false,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := repo.Get(tt.args.request)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricRepository_Set(t *testing.T) {
	c.ParseConfig()

	type args struct {
		request *domain.Metrics
	}
	req, repo := setup()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: req},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Set(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
