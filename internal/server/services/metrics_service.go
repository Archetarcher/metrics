package services

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"net/http"
)

type MetricsService struct {
	MetricRepository
}
type MetricRepository interface {
	GetAll() ([]models.Metrics, error)
	Get(request *models.Metrics) (*models.Metrics, *models.MetricError)
	Set(request *models.Metrics) error
}

func (s *MetricsService) Update(request *models.Metrics) (*models.Metrics, *models.MetricError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, err
	}
	if response != nil && request.MType == models.CounterType {
		c := *request.Delta + *response.Delta
		request.Delta = &c
	}
	if err := s.Set(request); err != nil {
		return nil, &models.MetricError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	response, _ = s.Get(request)
	return response, nil
}
func (s *MetricsService) GetValue(request *models.Metrics) (*models.Metrics, *models.MetricError) {

	response, err := s.Get(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, &models.MetricError{
			Code: http.StatusNotFound,
			Text: "value not found",
		}
	}

	return response, nil
}
func (s *MetricsService) GetAllValues() (string, *models.MetricError) {

	response, err := s.GetAll()

	if err != nil {
		return "", &models.MetricError{
			Code: http.StatusInternalServerError,
			Text: err.Error(),
		}
	}
	if response == nil {
		return "", &models.MetricError{
			Code: http.StatusNotFound,
			Text: "value not found",
		}
	}
	page := "<table><tr><th>Name</th><th>Value</th></tr>"

	for _, val := range response {
		v := models.GetStringValue(&val)
		page += "<tr><td>" + val.ID + "</td>" + "<td>" + v + "</td></tr>"
	}

	page += "</table>"
	return page, nil
}
