package store

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"strconv"
	"sync"
)

type MemStorage struct {
	mux  sync.Mutex
	data map[string]string
}

func NewStorage() *MemStorage {
	return &MemStorage{
		mux:  sync.Mutex{},
		data: make(map[string]string),
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
	if request.Type == domain.GaugeType {
		gaugeValue := request.Value

		s.data[getName(request)] = strconv.FormatFloat(gaugeValue, 'f', 3, 64)
		return nil
	}

	counterValue := int64(request.Value)
	s.data[getName(request)] = strconv.FormatInt(counterValue, 10)
	return nil
}

func getName(request *domain.MetricRequest) string {
	return fmt.Sprintf("%s_%s", request.Name, request.Type)
}
