package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackingService_FetchMemory(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchMemory()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTrackingService_FetchRuntime(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		counter int64
	}{
		{
			name:    "positive test #1",
			wantErr: false,
			counter: int64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchRuntime(tt.counter)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
