package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

type MetricRepository struct {
	Storage *store.MemStorage
}

func (r *MetricRepository) Get(request *domain.UpdateRequest) (*domain.MetricResponse, error) {
	return r.Storage.GetValue(request)
}
func (r *MetricRepository) Set(request *domain.UpdateRequest) error {
	return r.Storage.SetValue(request)
}
