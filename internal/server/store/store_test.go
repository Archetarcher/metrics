package store

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name string
		conf Config
		want string
	}{
		{
			name: "positive test #1",
			conf: Config{Memory: &memory.Config{}, Pgx: &pgx.Config{DatabaseDsn: ""}},
		},
		{
			name: "positive test #2",
			conf: Config{Pgx: &pgx.Config{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../migrations"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStore(context.Background(), tt.conf)
			assert.NotNil(t, s)
			assert.Nil(t, err)
		})
	}
}

func TestRetry(t *testing.T) {

	tests := []struct {
		name string
		conf Config
		want string
	}{
		{
			name: "positive test #1",
			conf: Config{Memory: &memory.Config{}, Pgx: &pgx.Config{DatabaseDsn: ""}},
		},
		{
			name: "positive test #2",
			conf: Config{Pgx: &pgx.Config{DatabaseDsn: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable", MigrationsPath: "../migrations"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Retry(context.Background(), &domain.MetricsError{}, 2, 3, tt.conf)
			assert.Nil(t, s)
			assert.NotNil(t, err)
		})
	}

}