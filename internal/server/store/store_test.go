package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMemStorage_GetValue(t *testing.T) {

	type args struct {
		request *models.Metrics
	}
	i := int64(1)
	tests := []struct {
		name    string
		args    args
		res     *models.Metrics
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				&models.Metrics{
					MType: "counter",
					ID:    "countervalue",
					Delta: &i,
				},
			},
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
		request *models.Metrics
	}
	i := int64(1)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				&models.Metrics{
					MType: "counter",
					ID:    "countervalue",
					Delta: &i,
				},
			},
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
				data: make(map[string]models.Metrics),
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
		request *models.Metrics
	}
	i := int64(1)

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{request: &models.Metrics{
				MType: "counter",
				ID:    "counterValue",
				Delta: &i,
			}},
			want: "counterValue_counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getName(tt.args.request))

		})
	}
}
