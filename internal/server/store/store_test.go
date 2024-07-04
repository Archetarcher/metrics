package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func setup() *domain.Metrics {
	i := int64(1)

	req := &domain.Metrics{
		MType: "counter",
		ID:    "counterValue",
		Delta: &i,
	}
	return req
}
func TestMemStorage_GetValue(t *testing.T) {

	type args struct {
		request *domain.Metrics
	}

	req := setup()
	tests := []struct {
		name    string
		args    args
		res     *domain.Metrics
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: req},
			wantErr: false,
			res:     nil,
		},
	}

	store := NewStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := store.GetValue(tt.args.request)

			assert.Equal(t, tt.res, res)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestMemStorage_SetValue(t *testing.T) {
	type args struct {
		request *domain.Metrics
	}
	req := setup()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: req},
			wantErr: false,
		},
	}

	store := NewStorage()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetValue(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "positive test #1",
			want: &MemStorage{
				mux:  sync.Mutex{},
				data: make(map[string]domain.Metrics),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewStorage())
		})
	}
}

func Test_getName(t *testing.T) {
	type args struct {
		request *domain.Metrics
	}
	req := setup()

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{request: req},
			want: "counterValue_counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getName(tt.args.request))

		})
	}
}
