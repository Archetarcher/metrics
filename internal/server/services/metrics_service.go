package services

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/utils"
	"net/http"
	"slices"
)

type MetricsService struct {
	repo MetricRepository
}

type MetricRepository interface {
	GetAllIn(keys []string) ([]domain.Metrics, *domain.MetricsError)
	GetAll() ([]domain.Metrics, *domain.MetricsError)
	Get(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	Set(request *domain.Metrics) *domain.MetricsError
	SetAll(request *[]domain.Metrics) *domain.MetricsError
}

func NewMetricsService(repo MetricRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) Updates(request *[]domain.Metrics) (*[]domain.Metrics, *domain.MetricsError) {
	keys := make([]string, len(*request))

	for _, m := range *request {
		if !slices.Contains(keys, getKey(m)) {
			keys = append(keys, getKey(m))
		}
	}

	metricsByKey, err := s.repo.GetAllIn(keys)
	if err != nil {
		return nil, err
	}

	existingKeys := make(map[string]int64)
	for key, m := range *request {
		if m.MType == domain.CounterType {
			if existingKeys[getKey(m)] != 0 {
				c := *m.Delta + existingKeys[getKey(m)]
				(*request)[key].Delta = &c
				continue
			}
			existingKeys[getKey(m)] = *m.Delta
		}

	}

	fmt.Println("metricsByKey")
	fmt.Println(metricsByKey)
	for _, mbk := range metricsByKey {
		for key, m := range *request {
			fmt.Println(getKey(m))
			fmt.Println(getKey(mbk))
			if getKey(m) == getKey(mbk) && m.MType == domain.CounterType {
				c := *m.Delta + *mbk.Delta
				fmt.Println("cccccc")
				fmt.Println(c)
				fmt.Println(getKey(m))
				(*request)[key].Delta = &c
			}
		}
	}

	if err := s.repo.SetAll(request); err != nil {
		return nil, err
	}
	response, _ := s.repo.GetAllIn(keys)

	return &response, nil
}
func (s *MetricsService) Update(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {
	response, err := s.repo.Get(request)
	if err != nil {
		return nil, err
	}
	if response != nil && request.MType == domain.CounterType {
		c := *request.Delta + *response.Delta
		request.Delta = &c
	}
	if err := s.repo.Set(request); err != nil {
		return nil, err
	}
	response, _ = s.repo.Get(request)
	return response, nil
}
func (s *MetricsService) GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {

	response, err := s.repo.Get(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, handleError(http.StatusNotFound, "value not found")
	}

	return response, nil
}
func (s *MetricsService) GetAllValues() (string, *domain.MetricsError) {
	response, err := s.repo.GetAll()

	if err != nil {
		return "", err
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
func getKey(request domain.Metrics) string {
	return fmt.Sprintf("%s_%s", request.ID, request.MType)
}
func handleError(code int, err string) *domain.MetricsError {
	return &domain.MetricsError{
		Code: code,
		Text: err,
	}
}
