package memory

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

const emptyParam = ""

// Store is a struct for in memory storage, keeps sync.Mutex and metrics map
type Store struct {
	data map[string]domain.Metrics
	mux  sync.Mutex
}

// NewStore creates new storage, restores data from file
func NewStore(ctx context.Context, config *config.AppConfig) (*Store, *domain.MetricsError) {
	logger.Log.Info("starting memory connection")

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
		storeInterval := time.Duration(config.StoreInterval) * time.Second

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

// RetryConnection retries connection to storage, not implemented for in memory storage
func RetryConnection(ctx context.Context, error *domain.MetricsError, interval int, try int, config *config.AppConfig) (*Store, *domain.MetricsError) {
	return nil, &domain.MetricsError{}
}

// CheckConnection checks connection to storage, not implemented for in memory storage
func (s *Store) CheckConnection(ctx context.Context) *domain.MetricsError {
	return nil
}

// Close closes connection to storage, not implemented for in memory storage
func (s *Store) Close() {

}

// GetValuesIn fetches metrics by keys in slices
func (s *Store) GetValuesIn(ctx context.Context, keys []string) ([]domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var metrics []domain.Metrics

	for _, k := range keys {
		metrics = append(metrics, s.data[k])
	}

	return metrics, nil
}

// GetValues fetches all metrics
func (s *Store) GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var res []domain.Metrics

	for _, v := range s.data {
		res = append(res, v)
	}
	return res, nil
}

// GetValue fetches metric by ID and MType from domain.Metrics
func (s *Store) GetValue(ctx context.Context, request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.data[getKey(*request)]
	if !ok {
		return nil, nil
	}
	return &res, nil
}

// SetValue sets metric data by key
func (s *Store) SetValue(ctx context.Context, request *domain.Metrics) *domain.MetricsError {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data[getKey(*request)] = *request
	return nil
}

// SetValues sets batch of metrics data by key
func (s *Store) SetValues(ctx context.Context, request []domain.Metrics) *domain.MetricsError {

	for _, v := range request {
		s.data[getKey(v)] = v
	}
	return nil
}

// Save saves metrics data to file
func (s *Store) Save(config *config.AppConfig) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if config.FileStoragePath == emptyParam {
		return nil
	}

	data, err := json.MarshalIndent(s.data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.FileStoragePath, data, 0666)
}

// Load loads metrics data from file
func (s *Store) Load(config *config.AppConfig) error {
	if config.FileStoragePath == emptyParam {
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
func getKey(request domain.Metrics) string {
	return request.ID + "_" + request.MType
}
