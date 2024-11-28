package utils

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"testing"
)

func TestGetStringValue(t *testing.T) {
	type args struct {
		result *domain.Metrics
	}
	c := int64(10)
	g := float64(10.12)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{&domain.Metrics{
				Delta: &c,
				Value: nil,
				ID:    "counter_value",
				MType: "counter",
			}},
			want: "10",
		},
		{
			name: "positive test #2",
			args: args{&domain.Metrics{
				Delta: nil,
				Value: &g,
				ID:    "gauge_value",
				MType: "gauge",
			}},
			want: "10.120",
		},
		{
			name: "positive test #3",
			args: args{&domain.Metrics{
				Delta: nil,
				Value: &g,
				ID:    "gauge_value",
				MType: "",
			}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStringValue(tt.args.result); got != tt.want {
				t.Errorf("GetStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
