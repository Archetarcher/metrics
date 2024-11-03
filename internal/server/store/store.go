package store

import (
	"context"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

const emptyParam = ""

// Store is interface that describes interaction with storage layer
type Store interface {
	GetValuesIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	SetValue(request *domain.Metrics, ctx context.Context) *domain.MetricsError
	SetValues(request []domain.Metrics, ctx context.Context) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
	Close()
}

// NewStore additional function to initiate Store instance according to factory pattern
func NewStore(conf Config, ctx context.Context) (Store, *domain.MetricsError) {

	if conf.Pgx.DatabaseDsn != emptyParam {
		return pgx.NewStore(conf.Pgx, ctx)
	}

	return memory.NewStore(conf.Memory, ctx)

}

// Retry retries connection to storage
func Retry(error *domain.MetricsError, interval int, try int, conf Config, ctx context.Context) (Store, *domain.MetricsError) {
	if conf.Pgx.DatabaseDsn != emptyParam {
		return pgx.RetryConnection(error, interval, try, conf.Pgx, ctx)
	}

	return memory.RetryConnection(error, interval, try, conf.Memory, ctx)
}
