package pgx

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"net/http"
	"sync"
	"time"
)

type Store struct {
	mux sync.Mutex
	db  *sql.DB
}

func NewStore(config *Config) (*Store, *domain.MetricsError) {

	db, err := sql.Open("pgx", config.DatabaseDsn)

	if err != nil {
		return nil, &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}

	storage := &Store{
		mux: sync.Mutex{},
		db:  db,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {

		return nil, &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}

	//if config.Restore {
	//	err := storage.Load(config)
	//	if err != nil {
	//		logger.Log.Info("failed to load metrics from file", zap.Error(err))
	//	}
	//}
	//
	//var wg sync.WaitGroup
	//wg.Add(1)
	//
	//go func() {
	//	defer wg.Done()
	//	var storeInterval = time.Duration(config.StoreInterval) * time.Second
	//
	//	for {
	//		err := storage.Save(config)
	//		if err != nil {
	//			logger.Log.Info("failed to store metrics to file", zap.Error(err))
	//		}
	//		time.Sleep(storeInterval)
	//	}
	//}()

	return storage, nil
}

func (s *Store) CheckConnection() *domain.MetricsError {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {

		return &domain.MetricsError{
			Text: err.Error(),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func (s *Store) GetValues() ([]domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	//var res []domain.Metrics
	//
	//for _, value := range s.data {
	//	res = append(res, value)
	//}
	return nil, nil
}
func (s *Store) GetValue(request *domain.Metrics) (*domain.Metrics, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	//res, ok := s.data[getName(request)]
	//if !ok {
	//	return nil, nil
	//}
	return nil, nil
}
func (s *Store) SetValue(request *domain.Metrics) error {
	//s.data[getName(request)] = *request
	return nil
}

func getName(request *domain.Metrics) string {
	return fmt.Sprintf("%s_%s", request.ID, request.MType)
}

func (s *Store) Save(config *Config) error {

	//if config.FileStoragePath == domain.EmptyParam {
	//	return nil
	//}
	//
	//data, err := json.MarshalIndent(s.data, "", "   ")
	//if err != nil {
	//	return err
	//}
	return nil
}
func (s *Store) Load(config *Config) error {
	//if config.FileStoragePath == domain.EmptyParam {
	//	return nil
	//}
	//
	//data, err := os.ReadFile(config.FileStoragePath)
	//
	//if errors.Is(err, os.ErrNotExist) {
	//	return nil
	//}
	//
	//if err != nil {
	//	return err
	//}
	//
	//var metrics map[string]domain.Metrics
	//if err := json.Unmarshal(data, &metrics); err != nil {
	//	return err
	//}
	//
	//s.data = metrics
	return nil
}
