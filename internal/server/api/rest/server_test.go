package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
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
			h := NewMetricsServer(&services.MetricsService{}, tt.fields.conf)

			hErr := h.Run()
			assert.Equal(t, tt.wantErr, hErr != nil, hErr)

		})
	}
}
