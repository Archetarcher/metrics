package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *AppConfig
	}{
		{
			name: "positive test 1",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig()
			if got == nil {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
