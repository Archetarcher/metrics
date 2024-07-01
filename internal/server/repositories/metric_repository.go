package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

type MetricRepository struct {
	Storage *store.MemStorage
}

func (r *MetricRepository) GetAll() ([]domain.MetricResponse, error) {
	return r.Storage.GetValues()
}
func (r *MetricRepository) Get(request *domain.MetricRequest) (*domain.MetricResponse, error) {
	return r.Storage.GetValue(request)
}
func (r *MetricRepository) Set(request *domain.MetricRequest) error {
	return r.Storage.SetValue(request)
}
