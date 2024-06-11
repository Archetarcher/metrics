package store

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"sync"
)

type MemStorage struct {
	mux  *sync.Mutex
	data map[string]float64
}

func NewStorage() *MemStorage {
	return &MemStorage{
		mux:  &sync.Mutex{},
		data: make(map[string]float64),
	}
}

func (s *MemStorage) GetValues() ([]domain.MetricResponse, error) {
	var res []domain.MetricResponse

	for name, value := range s.data {
		res = append(res, domain.MetricResponse{
			Name:  name,
			Value: value,
		})
	}
	return res, nil
}
func (s *MemStorage) GetValue(request *domain.MetricRequest) (*domain.MetricResponse, error) {
	res, ok := s.data[getName(request)]
	if !ok {
		return nil, nil
	}
	return &domain.MetricResponse{
		Name:  getName(request),
		Value: res,
	}, nil
}

func (s *MemStorage) SetValue(request *domain.MetricRequest) error {
	s.data[getName(request)] = request.Value
	return nil
}

func getName(request *domain.MetricRequest) string {
	return fmt.Sprintf("%s_%s", request.Name, request.Type)
}
