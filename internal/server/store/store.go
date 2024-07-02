package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"github.com/Archetarcher/metrics.git/internal/server/models"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type MemStorage struct {
	mux  sync.Mutex
	data map[string]models.Metrics
}

func NewStorage() *MemStorage {
	storage := &MemStorage{
		mux:  sync.Mutex{},
		data: make(map[string]models.Metrics),
	}

	if models.Restore {
		err := storage.Load()
		if err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var storeInterval = time.Duration(models.StoreInterval) * time.Second

		for {
			err := storage.Save()
			if err != nil {
				logger.Log.Info("failed to store metrics to file", zap.Error(err))
			}
			time.Sleep(storeInterval)
		}
	}()

	return storage
}

func (s *MemStorage) GetValues() ([]models.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var res []models.Metrics

	for _, value := range s.data {
		res = append(res, value)
	}
	return res, nil
}
func (s *MemStorage) GetValue(request *models.Metrics) (*models.Metrics, *models.MetricError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.data[getName(request)]
	if !ok {
		return nil, nil
	}
	return &res, nil
}
func (s *MemStorage) SetValue(request *models.Metrics) error {
	s.data[getName(request)] = *request
	return nil
}

func getName(request *models.Metrics) string {
	return fmt.Sprintf("%s_%s", request.ID, request.MType)
}

func (s *MemStorage) Save() error {

	if models.FileStoragePath == models.EmptyParam {
		return nil
	}

	data, err := json.MarshalIndent(s.data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(models.FileStoragePath, data, 0666)
}
func (s *MemStorage) Load() error {
	if models.FileStoragePath == models.EmptyParam {
		return nil
	}

	data, err := os.ReadFile(models.FileStoragePath)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	var metrics map[string]models.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	s.data = metrics
	return nil
}
