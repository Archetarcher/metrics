package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
)

func TestTrackingService_Fetch(t *testing.T) {
	type args struct {
		metrics         domain.MetricsData
		counterInterval int64
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				counterInterval: int64(1),
				metrics:         domain.MetricsData{},
			},
			wantErr: false,
		},
	}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.FetchRuntime(tt.args.counterInterval)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTrackingService_Send(t *testing.T) {
	c := config.AppConfig{}
	c.ParseConfig()
	type args struct {
		request []domain.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{}
	service := &TrackingService{Config: &c}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Send(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
