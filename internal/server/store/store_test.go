package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_GetValue(t *testing.T) {
	type fields struct {
		mux  *sync.Mutex
		data map[string]float64
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
			s := &MemStorage{
				mux:  tt.fields.mux,
				data: tt.fields.data,
			}
			got, err := s.GetValue(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_SetValue(t *testing.T) {
	type fields struct {
		mux  *sync.Mutex
		data map[string]float64
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
			s := &MemStorage{
				mux:  tt.fields.mux,
				data: tt.fields.data,
			}
			if err := s.SetValue(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getName(t *testing.T) {
	type args struct {
		request *domain.UpdateRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getName(tt.args.request); got != tt.want {
				t.Errorf("getName() = %v, want %v", got, tt.want)
			}
		})
	}
}
