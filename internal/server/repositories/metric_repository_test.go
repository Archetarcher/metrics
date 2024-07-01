package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricRepository_Get(t *testing.T) {

	type args struct {
		request *domain.MetricRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.MetricResponse
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				&domain.MetricRequest{
					Type:  "counter",
					Name:  "countervalue",
					Value: 1,
				},
			},
			wantErr: false,
			want:    nil,
		},
	}
	repo := &MetricRepository{Storage: store.NewStorage()}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := repo.Get(tt.args.request)

			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMetricRepository_Set(t *testing.T) {

	type args struct {
		request *domain.MetricRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				&domain.MetricRequest{
					Type:  "counter",
					Name:  "countervalue",
					Value: 1,
				},
			},
			wantErr: false,
		},
	}
	repo := &MetricRepository{Storage: store.NewStorage()}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Set(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
