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

func (s *MemStorage) GetValue(request *domain.UpdateRequest) (*domain.MetricResponse, error) {
	res, ok := s.data[getName(request)]
	if !ok {
		return nil, nil
	}
	return &domain.MetricResponse{
		Name:  getName(request),
		Value: res,
	}, nil

}
func (s *MemStorage) SetValue(request *domain.UpdateRequest) error {
	s.data[getName(request)] = request.Value
	return nil
}

func getName(request *domain.UpdateRequest) string {
	return fmt.Sprintf("%s_%s", request.Name, request.Type)
}
