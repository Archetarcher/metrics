package pgx

import (
	"context"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		c.c = &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migrations"}

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
	type args struct {
		request *config.AppConfig
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{request: &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migrations"}},
			wantErr: false,
		},
		{
			name:    "negative test #2",
			args:    args{request: &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migration"}},
			wantErr: true,
		},
		{
			name:    "negative test #3",
			args:    args{request: &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikums?sslmode=disable", MigrationsPath: "../../migrations"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStore(context.Background(), tt.args.request)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}

}

func TestRetryConnection(t *testing.T) {
	type args struct {
		config   *config.AppConfig
		err      *domain.MetricsError
		try      int
		interval int
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				config:   &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migrations"},
				try:      3,
				interval: 3,
				err: &domain.MetricsError{
					Err: &pgconn.PgError{Code: pgerrcode.ConnectionException},
				},
			},
			wantErr: false,
		},
		{
			name: "negative test, tries count exceeded #2",
			args: args{
				config:   &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migrations"},
				try:      0,
				interval: 3,
				err:      &domain.MetricsError{},
			},
			wantErr: true,
		},
		{
			name: "negative test, not connection error #3",
			args: args{
				config:   &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../../migrations"},
				try:      3,
				interval: 3,
				err:      &domain.MetricsError{},
			},
			wantErr: true,
		},
		{
			name: "negative test, incorrect connection string #4",
			args: args{
				config:   &config.AppConfig{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikums?sslmode=disable", MigrationsPath: "../../migrations"},
				try:      3,
				interval: 3,
				err: &domain.MetricsError{
					Err: &pgconn.PgError{Code: pgerrcode.ConnectionException},
				},
			},
			wantErr: true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RetryConnection(ctx, tt.args.err, tt.args.interval, tt.args.try, tt.args.config)
			fmt.Println(err)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
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

	type args struct {
		arr []string
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
			args:    args{arr: []string{values[0].ID + values[0].MType}},
		},
		{
			name:    "negative test #2",
			wantErr: true,
			args:    args{arr: []string{}},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conf.store.GetValuesIn(ctx, tt.args.arr)
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

func Test_runMigrations(t *testing.T) {
	require.Nil(t, conf.err, "failed to init store", conf.store, conf.err)

	type args struct {
		conf *config.AppConfig
	}

	tests := []struct {
		want    *domain.Metrics
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			args:    args{conf: conf.c},
			wantErr: false,
		},
		{
			name:    "negative test #2",
			args:    args{conf: &config.AppConfig{DatabaseDsn: conf.c.DatabaseDsn, MigrationsPath: ""}},
			wantErr: true,
		},
		{
			name:    "negative test #3",
			args:    args{conf: &config.AppConfig{MigrationsPath: conf.c.MigrationsPath, DatabaseDsn: "postgre://postgres:postgres@localhost:5432/praktikum?sslmode=disable"}},
			wantErr: true,
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cErr := runMigrations(ctx, tt.args.conf)
			assert.Equal(t, tt.wantErr, cErr != nil)
		})
	}
}
