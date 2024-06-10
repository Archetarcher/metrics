package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"net/http"
)

type MetricsService struct {
	MetricRepositoryInterface
}
type MetricRepositoryInterface interface {
	Get(request *domain.UpdateRequest) (*domain.MetricResponse, error)
	Set(request *domain.UpdateRequest) error
}

func (s *MetricsService) Update(request *domain.UpdateRequest) (*domain.MetricResponse, *domain.ApplicationError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response != nil && request.Type == domain.CounterType {
		request.Value += response.Value
	}
	if err := s.Set(request); err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	response, _ = s.Get(request)
	fmt.Println(response)
	return nil, nil
}
