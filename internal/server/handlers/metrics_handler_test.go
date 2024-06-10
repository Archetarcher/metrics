package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"net/http"
	"reflect"
	"testing"
)

func TestMetricsHandler_UpdateMetrics(t *testing.T) {
	type fields struct {
		MetricsServiceInterface MetricsServiceInterface
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandler{
				MetricsServiceInterface: tt.fields.MetricsServiceInterface,
			}
			h.UpdateMetrics(tt.args.w, tt.args.r)
		})
	}
}

func Test_validateRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name  string
		args  args
		want  *domain.UpdateRequest
		want1 *domain.ApplicationError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateRequest(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateRequest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("validateRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
