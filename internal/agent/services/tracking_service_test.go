package services

import (
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrackingService_Fetch(t *testing.T) {
	type args struct {
		counterInterval int64
		metrics         models.MetricsData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				counterInterval: int64(1),
				metrics:         models.MetricsData{},
			},
			wantErr: false,
		},
	}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Fetch(tt.args.counterInterval, &tt.args.metrics)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTrackingService_Send(t *testing.T) {
	config.ParseConfig()
	type args struct {
		request *models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Send(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
