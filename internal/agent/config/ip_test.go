package config

import (
	"net"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	tests := []struct {
		name string
		want net.IP
	}{
		{
			name: "positive test",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLocalIP()
			if got == nil {
				t.Errorf("GetLocalIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
