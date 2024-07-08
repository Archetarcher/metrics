package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

type MetricRepository struct {
	Storage *store.MemStorage
}

func (r *MetricRepository) GetAll() ([]domain.Metrics, error) {
	return r.Storage.GetValues()
}
func (r *MetricRepository) Get(request *domain.Metrics) (*domain.Metrics, error) {
	return r.Storage.GetValue(request)
}
func (r *MetricRepository) Set(request *domain.Metrics) error {
	return r.Storage.SetValue(request)
}
