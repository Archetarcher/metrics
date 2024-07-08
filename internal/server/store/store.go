package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type MemStorage struct {
	mux    sync.Mutex
	data   map[string]domain.Metrics
	Config *config.AppConfig
}

func NewStorage(config *config.AppConfig) *MemStorage {
	storage := &MemStorage{
		mux:    sync.Mutex{},
		data:   make(map[string]domain.Metrics),
		Config: config,
	}

	if config.Restore {
		err := storage.Load()
		if err != nil {
			logger.Log.Info("failed to load metrics from file", zap.Error(err))
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var storeInterval = time.Duration(config.StoreInterval) * time.Second

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

func (s *MemStorage) GetValues() ([]domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var res []domain.Metrics

	for _, value := range s.data {
		res = append(res, value)
	}
	return res, nil
}
func (s *MemStorage) GetValue(request *domain.Metrics) (*domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.data[getName(request)]
	if !ok {
		return nil, nil
	}
	return &res, nil
}
func (s *MemStorage) SetValue(request *domain.Metrics) error {
	s.data[getName(request)] = *request
	return nil
}

func getName(request *domain.Metrics) string {
	return fmt.Sprintf("%s_%s", request.ID, request.MType)
}

func (s *MemStorage) Save() error {

	if s.Config.FileStoragePath == domain.EmptyParam {
		return nil
	}

	data, err := json.MarshalIndent(s.data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Config.FileStoragePath, data, 0666)
}
func (s *MemStorage) Load() error {
	if s.Config.FileStoragePath == domain.EmptyParam {
		return nil
	}

	data, err := os.ReadFile(s.Config.FileStoragePath)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	var metrics map[string]domain.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}

	s.data = metrics
	return nil
}
