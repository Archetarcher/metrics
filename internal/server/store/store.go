package store

import (
	"context"
	"github.com/Archetarcher/metrics.git/internal/server/config"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

const emptyParam = ""

// Store is interface that describes interaction with storage layer
type Store interface {
	GetValuesIn(ctx context.Context, keys []string) ([]domain.Metrics, *domain.MetricsError)
	GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValue(ctx context.Context, request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	SetValue(ctx context.Context, request *domain.Metrics) *domain.MetricsError
	SetValues(ctx context.Context, request []domain.Metrics) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
}

// NewStore additional function to initiate Store instance according to factory pattern
func NewStore(ctx context.Context, conf *config.AppConfig) (Store, *domain.MetricsError) {

	if conf.DatabaseDsn != emptyParam {
		return pgx.NewStore(ctx, conf)
	}

	return memory.NewStore(ctx, conf)

}

// Retry retries connection to storage
func Retry(ctx context.Context, error *domain.MetricsError, interval int, try int, conf *config.AppConfig) (Store, *domain.MetricsError) {
	if conf.DatabaseDsn != emptyParam {
		return pgx.RetryConnection(ctx, error, interval, try, conf)
	}

	return memory.RetryConnection(ctx, error, interval, try, conf)
}
