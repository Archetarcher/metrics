package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
)

type MetricsService struct {
	Storage *store.MemStorage
}
type metricsServiceInterface interface {
	Update(request *domain.UpdateRequest) (*domain.MetricResponse, *domain.ApplicationError)
}

func (s *MetricsService) Update(request *domain.UpdateRequest) (*domain.MetricResponse, *domain.ApplicationError) {

	response, err := s.Storage.GetValue(request)
	if err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response != nil && request.Type == domain.CounterType {
		request.Value += response.Value
	}
	if err := s.Storage.SetValue(request); err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	response, _ = s.Storage.GetValue(request)
	fmt.Println(response)
	return nil, nil
}
