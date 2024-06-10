package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"reflect"
	"testing"
)

func TestMetricRepository_Get(t *testing.T) {
	type fields struct {
		Storage *store.MemStorage
	}
	type args struct {
		request *domain.UpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.MetricResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepository{
				Storage: tt.fields.Storage,
			}
			got, err := r.Get(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepository_Set(t *testing.T) {
	type fields struct {
		Storage *store.MemStorage
	}
	type args struct {
		request *domain.UpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepository{
				Storage: tt.fields.Storage,
			}
			if err := r.Set(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
