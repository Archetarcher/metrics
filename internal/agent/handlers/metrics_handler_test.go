package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"testing"
)

func TestNewMetricsHandler(t *testing.T) {
	type args struct {
		conf     *config.AppConfig
		provider MetricsProvider
		service  MetricsService
	}
	tests := []struct {
		name string
		args args
		want *MetricsHandler
	}{
		{
			name: "positive test 1",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetricsHandler(tt.args.conf, tt.args.provider, tt.args.service)
			if got == nil {
				t.Errorf("NewMetricsHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsHandler_StartSession(t *testing.T) {
	type fields struct {
		service  MetricsService
		provider MetricsProvider
		config   *config.AppConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *domain.MetricsError
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandler{
				service:  tt.fields.service,
				provider: tt.fields.provider,
				config:   tt.fields.config,
			}
			got := h.StartSession()
			if got != nil {
				t.Errorf("StartSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsHandler_TrackMetrics(t *testing.T) {
	type fields struct {
		service  MetricsService
		provider MetricsProvider
		config   *config.AppConfig
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandler{
				service:  tt.fields.service,
				provider: tt.fields.provider,
				config:   tt.fields.config,
			}
			h.TrackMetrics()
		})
	}
}
