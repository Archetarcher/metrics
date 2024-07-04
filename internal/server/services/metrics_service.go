package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/utils"
	"net/http"
)

type MetricsService struct {
	MetricRepository
}
type MetricRepository interface {
	GetAll() ([]domain.Metrics, error)
	Get(request *domain.Metrics) (*domain.Metrics, error)
	Set(request *domain.Metrics) error
}

func (s *MetricsService) Update(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, handleError(http.StatusInternalServerError, err.Error())
	}
	if response != nil && request.MType == domain.CounterType {
		c := *request.Delta + *response.Delta
		request.Delta = &c
	}
	if err := s.Set(request); err != nil {
		return nil, handleError(http.StatusNotFound, err.Error())
	}
	response, _ = s.Get(request)
	return response, nil
}
func (s *MetricsService) GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, handleError(http.StatusInternalServerError, err.Error())
	}
	if response == nil {
		return nil, handleError(http.StatusNotFound, "value not found")
	}

	return response, nil
}
func (s *MetricsService) GetAllValues() (string, *domain.MetricsError) {

	response, err := s.GetAll()

	if err != nil {
		return "", handleError(http.StatusInternalServerError, err.Error())
	}
	if response == nil {
		return "", handleError(http.StatusNotFound, "value not found")
	}
	page := "<table><tr><th>Name</th><th>Value</th></tr>"

	for _, val := range response {
		v := utils.GetStringValue(&val)
		page += "<tr><td>" + val.ID + "</td>" + "<td>" + v + "</td></tr>"
	}

	page += "</table>"
	return page, nil
}
func handleError(code int, err string) *domain.MetricsError {
	return &domain.MetricsError{
		Code: code,
		Text: err,
	}
}
