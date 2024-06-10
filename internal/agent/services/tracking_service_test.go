package services

import (
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"reflect"
	"testing"
)

func TestTrackingService_Fetch(t *testing.T) {
	type args struct {
		counterInterval int
	}
	tests := []struct {
		name  string
		args  args
		want  []domain.MetricData
		want1 *domain.ApplicationError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TrackingService{}
			got, got1 := s.Fetch(tt.args.counterInterval)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Fetch() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTrackingService_Send(t *testing.T) {
	type args struct {
		request *domain.MetricData
	}
	tests := []struct {
		name  string
		args  args
		want  *domain.ServerResponse
		want1 *domain.ApplicationError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TrackingService{}
			got, got1 := s.Send(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Send() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Send() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
