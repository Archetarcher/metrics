package handlers

import (
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"sync"
	"testing"
)

func TestTrackingHandler_StartTracking(t *testing.T) {
	type fields struct {
		TrackingServiceInterface TrackingServiceInterface
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TrackingHandler{
				TrackingServiceInterface: tt.fields.TrackingServiceInterface,
			}
			h.StartTracking()
		})
	}
}

func Test_startPoll(t *testing.T) {
	type args struct {
		fetch   fetch
		metrics chan<- domain.MetricData
		wg      *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startPoll(tt.args.fetch, tt.args.metrics, tt.args.wg)
		})
	}
}

func Test_startReport(t *testing.T) {
	type args struct {
		send    send
		metrics <-chan domain.MetricData
		wg      *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startReport(tt.args.send, tt.args.metrics, tt.args.wg)
		})
	}
}
