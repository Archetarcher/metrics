package repositories

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

type MetricRepository struct {
	Storage *store.MemStorage
}

func (r *MetricRepository) GetAll() ([]models.Metrics, error) {
	return r.Storage.GetValues()
}
func (r *MetricRepository) Get(request *models.Metrics) (*models.Metrics, *models.MetricError) {
	return r.Storage.GetValue(request)
}
func (r *MetricRepository) Set(request *models.Metrics) error {
	return r.Storage.SetValue(request)
}
