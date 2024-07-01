package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/repositories"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsService_Update(t *testing.T) {

	type args struct {
		request *domain.MetricRequest
	}
	tests := []struct {
		name string
		args args
		res  *domain.MetricResponse
		err  *domain.MetricError
	}{
		{
			name: "positive test #1",
			args: args{request: &domain.MetricRequest{Type: "gauge", Name: "test", Value: 1}},
			res:  nil,
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
