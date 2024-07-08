package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var c = config.NewConfig()

func setup() (*domain.Metrics, *domain.Metrics) {
	i := float64(1)

	req := &domain.Metrics{MType: "gauge", ID: "test", Value: &i}
	res := &domain.Metrics{MType: "gauge", ID: "test", Value: &i}
	return req, res
}
func TestMetricsService_Update(t *testing.T) {
	c.ParseConfig()

	type args struct {
		request *domain.Metrics
	}
	req, res := setup()
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
	repo := &repositories.MetricRepository{Storage: store.NewStorage(c)}
	service := &MetricsService{MetricRepository: repo}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.Update(tt.args.request)
			assert.Equal(t, tt.res, res)
			assert.Equal(t, tt.err, err)

		})
	}
}
