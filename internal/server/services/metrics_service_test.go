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
	"go.uber.org/zap"
	"testing"
)

var c = config.NewConfig(store.Config{Memory: &memory.Config{Active: true}, Pgx: &pgx.Config{}})

func setup() (*domain.Metrics, *domain.Metrics, *MetricsService) {
	i := float64(1)

	req := &domain.Metrics{MType: "gauge", ID: "test", Value: &i}
	res := &domain.Metrics{MType: "gauge", ID: "test", Value: &i}

	ctx := context.Background()
	storage, err := store.NewStore(c.Store, ctx)
	if err != nil {
		logger.Log.Error("failed to init storage with error", zap.String("error", err.Text), zap.Int("code", err.Code))
	}

	repo := repositories.NewMetricsRepository(storage)
	service := NewMetricsService(repo)
	return req, res, service
}
func TestMetricsService_Update(t *testing.T) {
	c.ParseConfig()

	type args struct {
		request *domain.Metrics
	}
	req, res, service := setup()
	ctx := context.Background()

	tests := []struct {
		name string
		args args
		res  *domain.Metrics
		err  *domain.MetricsError
	}{
		{
			name: "positive test #1",
			args: args{request: req},
			res:  res,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.Update(tt.args.request, ctx)
			assert.Equal(t, tt.res, res)
			assert.Equal(t, tt.err, err)

		})
	}
}
