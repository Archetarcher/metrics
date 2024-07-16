package memory

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
func TestStore_GetValue(t *testing.T) {

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

	store := NewStore(&Config{StoreInterval: 300, FileStoragePath: "/tmp/metrics-pgx.json", Restore: false})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := store.GetValue(tt.args.request)

			assert.Equal(t, tt.res, res)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestStore_SetValue(t *testing.T) {

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

	store := NewStore(&Config{StoreInterval: 300, FileStoragePath: "/tmp/metrics-pgx.json", Restore: false})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetValue(tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestNewStore(t *testing.T) {
	store := NewStore(&Config{StoreInterval: 300, FileStoragePath: "/tmp/metrics-pgx.json", Restore: false})

	tests := []struct {
		name string
		want *Store
	}{
		{
			name: "positive test #1",
			want: &Store{
				mux:  sync.Mutex{},
				data: make(map[string]domain.Metrics),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, store)
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
