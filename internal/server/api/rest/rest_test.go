package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
	"reflect"
	"testing"
)

func TestAPI_Run(t *testing.T) {
	type fields struct {
		server *http.ServeMux
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := API{
				server: tt.fields.server,
			}
			if err := a.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAPI(t *testing.T) {
	type args struct {
		storage *store.MemStorage
	}
	tests := []struct {
		name string
		args args
		want API
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPI(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
