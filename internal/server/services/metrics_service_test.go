package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsService_Update(t *testing.T) {

	type args struct {
		request *models.Metrics
	}
	i := float64(1)
	tests := []struct {
		name string
		args args
		res  *models.Metrics
		err  *models.MetricError
	}{
		{
			name: "positive test #1",
			args: args{request: &models.Metrics{MType: "gauge", ID: "test", Value: &i}},
			res:  &models.Metrics{MType: "gauge", ID: "test", Value: &i},
			err:  nil,
		},
	}
	repo := &repositories.MetricRepository{Storage: store.NewStorage()}
	service := &MetricsService{MetricRepository: repo}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.Update(tt.args.request)
			assert.Equal(t, tt.res, res)
			assert.Equal(t, tt.err, err)

		})
	}
}
