package services

import (
	"context"
	"net/http"
	"slices"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/utils"
)

// MetricsService is a service struct for metrics, keeps implementation of MetricRepository interface
type MetricsService struct {
	repo MetricRepository
}

// MetricRepository is an interface that describes interaction with repository layer
type MetricRepository interface {
	GetAllIn(ctx context.Context, keys []string) ([]domain.Metrics, *domain.MetricsError)
	GetAll(ctx context.Context) ([]domain.Metrics, *domain.MetricsError)
	Get(ctx context.Context, request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	Set(ctx context.Context, request *domain.Metrics) *domain.MetricsError
	SetAll(ctx context.Context, request []domain.Metrics) *domain.MetricsError
	CheckConnection(ctx context.Context) *domain.MetricsError
}

// NewMetricsService creates MetricsService
func NewMetricsService(repo MetricRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

// CheckConnection checks connection to storage in repository
func (s *MetricsService) CheckConnection(ctx context.Context) *domain.MetricsError {
	return s.repo.CheckConnection(ctx)
}

// Updates creates or updates batch of metrics data
func (s *MetricsService) Updates(request []domain.Metrics, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	keys := make([]string, len(request))

	for _, m := range request {
		if !slices.Contains(keys, getKey(m)) {
			keys = append(keys, getKey(m))
		}
	}

	metricsByKey, err := s.repo.GetAllIn(ctx, keys)
	if err != nil {
		return nil, err
	}

	existingKeys := make(map[string]int64)
	for k, m := range request {
		if m.MType == domain.CounterType {
			if existingKeys[getKey(m)] != 0 {
				c := *m.Delta + existingKeys[getKey(m)]
				(request)[k].Delta = &c
				continue
			}
			existingKeys[getKey(m)] = *m.Delta
		}

	}

	for _, mbk := range metricsByKey {
		for k, m := range request {
			if getKey(m) == getKey(mbk) && m.MType == domain.CounterType {
				c := *m.Delta + *mbk.Delta
				(request)[k].Delta = &c
			}
		}
	}

	if rErr := s.repo.SetAll(ctx, request); rErr != nil {
		return nil, rErr
	}
	response, err := s.repo.GetAllIn(ctx, keys)

	if err != nil {
		return nil, err
	}
	return response, nil
}

// Update creates or updates metric data
func (s *MetricsService) Update(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	response, err := s.repo.Get(ctx, request)
	if err != nil {
		return nil, err
	}

	if response != nil && request.MType == domain.CounterType {
		c := *request.Delta + *response.Delta
		request.Delta = &c
	}

	if rErr := s.repo.Set(ctx, request); rErr != nil {
		return nil, rErr
	}

	response, err = s.repo.Get(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetValue fetches metric data by ID and MType in domain.Metrics
func (s *MetricsService) GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {

	response, err := s.repo.Get(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, handleError(http.StatusNotFound, "value not found")
	}

	return response, nil
}

// GetAllValues fetches all metrics data
func (s *MetricsService) GetAllValues(ctx context.Context) (string, *domain.MetricsError) {
	response, err := s.repo.GetAll(ctx)

	if err != nil {
		return "", err
	}
	if response == nil {
		return "", handleError(http.StatusNotFound, "value not found")
	}
	page := "<table><tr><th>Name</th><th>Value</th></tr>"

	for _, r := range response {
		v := utils.GetStringValue(&r)
		page += "<tr><td>" + r.ID + "</td>" + "<td>" + v + "</td></tr>"
	}

	page += "</table>"
	return page, nil
}

func getKey(request domain.Metrics) string {
	return request.ID + "_" + request.MType
}
func handleError(code int, err string) *domain.MetricsError {
	return &domain.MetricsError{
		Code: code,
		Text: err,
	}
}
