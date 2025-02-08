package memory

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

var conf TConfig

type TConfig struct {
	c     *config.AppConfig
	store *Store
	err   *domain.MetricsError
	once  sync.Once
}

func (c *TConfig) setConfig() {
	c.once.Do(func() {
		c.c = &config.AppConfig{}

		s, err := NewStore(context.Background(), c.c)

		c.store = s
		c.err = err
	})
}
func init() {
	conf.setConfig()
}

var (
	counter = int64(2896127014)
	gauge   = 0.31167763133187076
	values  = [8]domain.Metrics{
		{
			ID:    "counter_value",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_2",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_2",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_3",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_3",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
		{
			ID:    "counter_value_4",
			MType: "counter",
			Delta: &counter,
			Value: nil,
		},
		{
			ID:    "gauge_value_4",
			MType: "gauge",
			Delta: nil,
			Value: &gauge,
		},
	}
)

func TestNewStore(t *testing.T) {
	t.Run("positive test", func(t *testing.T) {
		s, err := NewStore(context.Background(), &config.AppConfig{})
		assert.NotNil(t, s)
		assert.Nil(t, err)
	})
}

func TestRetryConnection(t *testing.T) {
	t.Run("positive test", func(t *testing.T) {
		s, err := RetryConnection(context.Background(), &domain.MetricsError{}, 3, 3, &config.AppConfig{})
		assert.Nil(t, s)
		assert.NotNil(t, err)
	})
}

func TestStore_CheckConnection(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	t.Run("positive test", func(t *testing.T) {

		cErr := conf.store.CheckConnection(context.Background())
		assert.Nil(t, cErr)
	})
}

func TestStore_GetValue(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	type args struct {
		request *domain.Metrics
	}

	tests := []struct {
		want    *domain.Metrics
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: &values[0]},
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.store.GetValue(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestStore_GetValues(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	type args struct {
	}

	tests := []struct {
		want    *domain.Metrics
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.store.GetValues(ctx)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestStore_GetValuesIn(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	tests := []struct {
		want    *domain.Metrics
		args    []string
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			wantErr: false,
			args:    []string{values[0].ID + values[0].MType},
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.store.GetValuesIn(ctx, tt.args)

			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}

func TestStore_SetValue(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	type args struct {
		request *domain.Metrics
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: &values[0]},
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.store.SetValue(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestStore_SetValues(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	type args struct {
		request []domain.Metrics
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: []domain.Metrics{values[1], values[2]}},
			wantErr: false,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.store.SetValues(ctx, tt.args.request)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_getKey(t *testing.T) {
	type args struct {
		request domain.Metrics
	}
	c := int64(10)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{domain.Metrics{
				Delta: &c,
				Value: nil,
				ID:    "value",
				MType: "counter",
			}},
			want: "value_counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getKey(tt.args.request); got != tt.want {
				t.Errorf("getKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Save(t *testing.T) {
	tests := []struct {
		conf    *config.AppConfig
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			conf:    &config.AppConfig{FileStoragePath: ""},
			wantErr: false,
		},
		{
			name:    "positive test #2",
			conf:    &config.AppConfig{FileStoragePath: "../../../../metrics.json"},
			wantErr: false,
		},
		{
			name: "positive test #3",
			conf: &config.AppConfig{FileStoragePath: "../../../../m.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.store.Load(tt.conf)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestStore_Load(t *testing.T) {
	tests := []struct {
		conf    *config.AppConfig
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			conf:    &config.AppConfig{FileStoragePath: ""},
			wantErr: false,
		}, {
			name:    "negative test #2",
			conf:    &config.AppConfig{FileStoragePath: "../../../../metrics.jsn"},
			wantErr: false,
		},
		{
			name:    "positive test #3",
			conf:    &config.AppConfig{FileStoragePath: "../../../../metrics.json"},
			wantErr: false,
		},
		{
			name:    "positive test #3",
			conf:    &config.AppConfig{FileStoragePath: "../../../../metrics2.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conf.store.Load(tt.conf)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
