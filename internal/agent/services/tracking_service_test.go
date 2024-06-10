package services

import (
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrackingService_Fetch(t *testing.T) {
	type args struct {
		counterInterval int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				counterInterval: 1,
			},
			wantErr: false,
		},
	}
	service := &TrackingService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.Fetch(tt.args.counterInterval)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.NotNil(t, res)
		})
	}
}

func TestTrackingService_Send(t *testing.T) {
	type args struct {
		request *domain.MetricData
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
