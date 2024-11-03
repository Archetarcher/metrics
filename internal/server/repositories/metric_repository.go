package repositories

import (
	"context"
	"time"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

// MetricRepository is a repository struct for metrics, keeps store interfaces implementation
type MetricRepository struct {
	store Store
}

// Store is an interface that describes interaction with storage layer
type Store interface {
	GetValuesIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	SetValue(request *domain.Metrics, ctx context.Context) *domain.MetricsError
	SetValues(request []domain.Metrics, ctx context.Context) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
	Close()
}

// NewMetricsRepository creates MetricRepository
func NewMetricsRepository(store Store) *MetricRepository {
	return &MetricRepository{
		store: store,
	}
}

// CheckConnection checks connection to storage
func (r *MetricRepository) CheckConnection(ctx context.Context) *domain.MetricsError {
	return r.store.CheckConnection(ctx)
}

// GetAllIn fetches all metrics with keys equivalent to keys in slice
func (r *MetricRepository) GetAllIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.store.GetValuesIn(keys, ctx)
}

// GetAll fetches all metrics
func (r *MetricRepository) GetAll(ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.store.GetValues(ctx)
}

// Get fetches one metric by ID and MType in domain.Metrics
func (r *MetricRepository) Get(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.store.GetValue(request, ctx)
}

// Set creates or updates metric data in storage
func (r *MetricRepository) Set(request *domain.Metrics, ctx context.Context) *domain.MetricsError {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.store.SetValue(request, ctx)
}

// SetAll creates or updates batch of metrics data in storage
func (r *MetricRepository) SetAll(request []domain.Metrics, ctx context.Context) *domain.MetricsError {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.store.SetValues(request, ctx)
}
