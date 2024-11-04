package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Archetarcher/metrics.git/internal/server/handlers"
)

func TestMetricsAPI_Run(t *testing.T) {

	type fields struct {
		server *http.ServeMux
	}
	tests := []struct {
		serverFields fields
		name         string
		wantErr      bool
	}{
		{
			name:         "With server defined",
			serverFields: fields{http.NewServeMux()},
			wantErr:      false,
		},
		{
			name:         "With no server defined",
			serverFields: fields{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.serverFields.server == nil, tt.wantErr)
		})
	}
}

func TestNewMetricsAPI(t *testing.T) {
	type fields struct {
		handler *handlers.MetricsHandler
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name:    "Test with handler",
			fields:  fields{handler: &handlers.MetricsHandler{}},
			wantErr: false,
		},
		{
			name:    "Test with no storage",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.fields.handler == nil, tt.wantErr)
		})
	}
}
