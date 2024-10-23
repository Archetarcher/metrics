package memory

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type Store struct {
	mux  sync.Mutex
	data map[string]domain.Metrics
}

func NewStore(config *Config, ctx context.Context) (*Store, *domain.MetricsError) {
	storage := &Store{
		mux:  sync.Mutex{},
		data: make(map[string]domain.Metrics),
	}

	if config.Restore {
		err := storage.Load(config)
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
			err := storage.Save(config)
			if err != nil {
				logger.Log.Info("failed to store metrics to file", zap.Error(err))
			}
			time.Sleep(storeInterval)
		}
	}()

	return storage, nil
}
func RetryConnection(error *domain.MetricsError, interval int, try int, config *Config, ctx context.Context) (*Store, *domain.MetricsError) {
	return nil, nil
}

func (s *Store) CheckConnection(ctx context.Context) *domain.MetricsError {
	return nil
}

func (s *Store) Close() {

}
func (s *Store) GetValuesIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var metrics []domain.Metrics

	for _, key := range keys {
		metrics = append(metrics, s.data[key])
	}

	return metrics, nil
}
func (s *Store) GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var res []domain.Metrics

	for _, value := range s.data {
		res = append(res, value)
	}
	return res, nil
}
func (s *Store) GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.data[getName(*request)]
	if !ok {
		return nil, nil
	}
	return &res, nil
}
func (s *Store) SetValue(request *domain.Metrics, ctx context.Context) *domain.MetricsError {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data[getName(*request)] = *request
	return nil
}
func (s *Store) SetValues(request []domain.Metrics, ctx context.Context) *domain.MetricsError {
	for _, value := range request {
		s.data[getName(value)] = value
	}
	return nil
}

func getName(request domain.Metrics) string {
	return request.ID + "_" + request.MType
}

func (s *Store) Save(config *Config) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if config.FileStoragePath == domain.EmptyParam {
		return nil
	}

	data, err := json.MarshalIndent(s.data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.FileStoragePath, data, 0666)
}
func (s *Store) Load(config *Config) error {
	if config.FileStoragePath == domain.EmptyParam {
		return nil
	}

	data, err := os.ReadFile(config.FileStoragePath)

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

func handleError(text string, code int) *domain.MetricsError {
	return &domain.MetricsError{
		Text: text,
		Code: code,
	}
}
