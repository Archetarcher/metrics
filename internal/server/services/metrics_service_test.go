package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"reflect"
	"testing"
)

func TestMetricsService_Update(t *testing.T) {
	type fields struct {
		MetricRepositoryInterface MetricRepositoryInterface
	}
	type args struct {
		request *domain.UpdateRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.MetricResponse
		want1  *domain.ApplicationError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{
				MetricRepositoryInterface: tt.fields.MetricRepositoryInterface,
			}
			got, got1 := s.Update(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Update() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
