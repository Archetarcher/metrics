package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Archetarcher/metrics.git/internal/server/handlers"
)

func TestNewMetricsAPI(t *testing.T) {
	type fields struct {
		handler *handlers.MetricsHandler
		conf    *config.AppConfig
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name:    "Test with handler",
			fields:  fields{handler: &handlers.MetricsHandler{}, conf: &config.AppConfig{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMetricsAPI(tt.fields.handler, tt.fields.conf)
			assert.Nil(t, err)

		})
	}
}

func TestMetricsAPI_Run(t *testing.T) {

	type fields struct {
		conf *config.AppConfig
	}

	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name:    "Negative test #2",
			fields:  fields{conf: &config.AppConfig{RunAddr: "8080"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewMetricsAPI(&handlers.MetricsHandler{}, tt.fields.conf)
			assert.Nil(t, err)

			hErr := h.Run(tt.fields.conf)
			assert.Equal(t, tt.wantErr, hErr != nil, hErr)

		})
	}
}
