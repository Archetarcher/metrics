package repositories

import (
	"context"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

type MetricRepository struct {
	store Store
}

type Store interface {
	GetValuesIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	SetValue(request *domain.Metrics, ctx context.Context) *domain.MetricsError
	SetValues(request []domain.Metrics, ctx context.Context) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
	Close()
}

func NewMetricsRepository(store Store) *MetricRepository {
	return &MetricRepository{
		store: store,
	}
}

func (r *MetricRepository) CheckConnection(ctx context.Context) *domain.MetricsError {
	return r.store.CheckConnection(ctx)
}
func (r *MetricRepository) GetAllIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	return r.store.GetValuesIn(keys, ctx)
}
func (r *MetricRepository) GetAll(ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	return r.store.GetValues(ctx)
}
func (r *MetricRepository) Get(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	return r.store.GetValue(request, ctx)
}
func (r *MetricRepository) Set(request *domain.Metrics, ctx context.Context) *domain.MetricsError {
	return r.store.SetValue(request, ctx)
}
func (r *MetricRepository) SetAll(request []domain.Metrics, ctx context.Context) *domain.MetricsError {
	return r.store.SetValues(request, ctx)
}
