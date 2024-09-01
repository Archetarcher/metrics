package services

import (
	"context"
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
	GetAllIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	GetAll(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	Get(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError)
	Set(request *domain.Metrics, ctx context.Context) *domain.MetricsError
	SetAll(request []domain.Metrics, ctx context.Context) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
}

func NewMetricsService(repo MetricRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) CheckConnection(ctx context.Context) *domain.MetricsError {
	return s.repo.CheckConnection(ctx)
}
func (s *MetricsService) Updates(request []domain.Metrics, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	keys := make([]string, len(request))

	for _, m := range request {
		if !slices.Contains(keys, getKey(m)) {
			keys = append(keys, getKey(m))
		}
	}

	metricsByKey, err := s.repo.GetAllIn(keys, ctx)
	if err != nil {
		return nil, err
	}

	existingKeys := make(map[string]int64)
	for key, m := range request {
		if m.MType == domain.CounterType {
			if existingKeys[getKey(m)] != 0 {
				c := *m.Delta + existingKeys[getKey(m)]
				(request)[key].Delta = &c
				continue
			}
			existingKeys[getKey(m)] = *m.Delta
		}

	}

	for _, mbk := range metricsByKey {
		for key, m := range request {
			if getKey(m) == getKey(mbk) && m.MType == domain.CounterType {
				c := *m.Delta + *mbk.Delta
				(request)[key].Delta = &c
			}
		}
	}

	if err := s.repo.SetAll(request, ctx); err != nil {
		return nil, err
	}
	response, _ := s.repo.GetAllIn(keys, ctx)

	return response, nil
}
func (s *MetricsService) Update(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	response, err := s.repo.Get(request, ctx)
	if err != nil {
		return nil, err
	}

	if response != nil && request.MType == domain.CounterType {
		c := *request.Delta + *response.Delta
		request.Delta = &c
	}

	if err := s.repo.Set(request, ctx); err != nil {
		return nil, err
	}

	response, err = s.repo.Get(request, ctx)
	if err := s.repo.Set(request, ctx); err != nil {
		return nil, err
	}
	return response, nil
}
func (s *MetricsService) GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {

	response, err := s.repo.Get(request, ctx)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, handleError(http.StatusNotFound, "value not found")
	}

	return response, nil
}
func (s *MetricsService) GetAllValues(ctx context.Context) (string, *domain.MetricsError) {
	response, err := s.repo.GetAll(ctx)

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
