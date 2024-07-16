package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

type MetricRepository struct {
	store Store
}

type Store interface {
	GetValues() ([]domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	SetValue(request *domain.Metrics) *domain.MetricsError
	CheckConnection() *domain.MetricsError
	Close()
}

func NewMetricsRepository(store Store) *MetricRepository {
	return &MetricRepository{
		store: store,
	}
}

func (r *MetricRepository) GetAll() ([]domain.Metrics, *domain.MetricsError) {
	return r.store.GetValues()
}
func (r *MetricRepository) Get(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {
	return r.store.GetValue(request)
}
func (r *MetricRepository) Set(request *domain.Metrics) *domain.MetricsError {
	return r.store.SetValue(request)
}
