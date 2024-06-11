package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"net/http"
	"strconv"
)

type MetricsService struct {
	MetricRepositoryInterface
}
type MetricRepositoryInterface interface {
	GetAll() ([]domain.MetricResponse, error)
	Get(request *domain.MetricRequest) (*domain.MetricResponse, error)
	Set(request *domain.MetricRequest) error
}

func (s *MetricsService) Update(request *domain.MetricRequest) (*domain.MetricResponse, *domain.ApplicationError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response != nil && request.Type == domain.CounterType {
		i, err := strconv.ParseFloat(response.Value, 64)
		if err != nil {
			return nil, &domain.ApplicationError{
				Code: http.StatusInternalServerError,
				Text: err.Error(),
			}
		}
		request.Value += i
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
func (s *MetricsService) GetValue(request *domain.MetricRequest) (*domain.MetricResponse, *domain.ApplicationError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response == nil {
		return nil, &domain.ApplicationError{
			Code: http.StatusNotFound,
			Text: "value not found",
		}
	}

	return &domain.MetricResponse{
		Name:  response.Name,
		Value: response.Value,
	}, nil
}
func (s *MetricsService) GetAllValues() (string, *domain.ApplicationError) {

	response, err := s.GetAll()

	if err != nil {
		return "", &domain.ApplicationError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response == nil {
		return "", &domain.ApplicationError{
			Code: http.StatusNotFound,
			Text: "value not found",
		}
	}
	page := "<table><tr><th>Name</th><th>Value</th></tr>"

	for _, val := range response {
		page += "<tr><td>" + val.Name + "</td>" + "<td>" + val.Value + "</td></tr>"
	}

	page += "</table>"
	return page, nil
}
