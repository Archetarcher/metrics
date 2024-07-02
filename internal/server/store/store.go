package store

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"sync"
)

type MemStorage struct {
	mux  sync.Mutex
	data map[string]models.Metrics
}

func NewStorage() *MemStorage {
	return &MemStorage{
		mux:  sync.Mutex{},
		data: make(map[string]models.Metrics),
	}
}

func (s *MemStorage) GetValues() ([]models.Metrics, error) {
	var res []models.Metrics

	for _, value := range s.data {
		res = append(res, value)
	}
	return res, nil
}
func (s *MemStorage) GetValue(request *models.Metrics) (*models.Metrics, *models.MetricError) {
	res, ok := s.data[getName(request)]
	if !ok {
		return nil, nil
	}

	//if request.MType == models.GaugeType {
	//	v, err := strconv.ParseFloat(res, 64)
	//
	//	if err != nil {
	//		return nil, &models.MetricError{Code: http.StatusInternalServerError}
	//	}
	//	request.Value = &v
	//}
	//if request.MType == models.CounterType {
	//	v, err := strconv.ParseInt(res, 10, 64)
	//
	//	if err != nil {
	//		return nil, &models.MetricError{Code: http.StatusInternalServerError}
	//	}
	//	request.Delta = &v
	//}
	return &res, nil
}

func (s *MemStorage) SetValue(request *models.Metrics) error {
	//if request.MType == models.GaugeType {
	//	gaugeValue := request.Value
	//
	//	s.data[getName(request)] = strconv.FormatFloat(*gaugeValue, 'f', 3, 64)
	//	return nil
	//}
	//
	//counterValue := request.Delta
	//s.data[getName(request)] = strconv.FormatInt(*counterValue, 10)
	s.data[getName(request)] = *request
	return nil
}

func getName(request *models.Metrics) string {
	return fmt.Sprintf("%s_%s", request.ID, request.MType)
}
