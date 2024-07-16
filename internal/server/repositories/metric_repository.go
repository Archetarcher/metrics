package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

type MetricRepository struct {
	store Store
}

type Store interface {
	GetValues() ([]domain.Metrics, error)
	GetValue(request *domain.Metrics) (*domain.Metrics, error)
	SetValue(request *domain.Metrics) error
	CheckConnection() *domain.MetricsError
}

func NewMetricsRepository(store Store) *MetricRepository {
	return &MetricRepository{
		store: store,
	}

}

func (r *MetricRepository) GetAll() ([]domain.Metrics, error) {
	return r.store.GetValues()
}
func (r *MetricRepository) Get(request *domain.Metrics) (*domain.Metrics, error) {
	return r.store.GetValue(request)
}
func (r *MetricRepository) Set(request *domain.Metrics) error {
	return r.store.SetValue(request)
}
