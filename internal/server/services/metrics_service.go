package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
)

type metricsService struct {
	storage *store.MemStorage
}
type metricsServiceInterface interface {
	Update(request *domain.UpdateRequest) (*domain.MetricResponse, *domain.ApplicationError)
}

var (
	MetricsService metricsServiceInterface
)

func init() {
	MetricsService = &metricsService{
		store.NewStorage(),
	}
}

func (s *metricsService) Update(request *domain.UpdateRequest) (*domain.MetricResponse, *domain.ApplicationError) {

	response, err := s.storage.GetValue(request)
	if err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response != nil && request.Type == domain.CounterType {
		request.Value += response.Value
	}
	if err := s.storage.SetValue(request); err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	response, _ = s.storage.GetValue(request)
	fmt.Println(response)
	return nil, nil
}
