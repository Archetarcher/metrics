package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(config *Config) (*Store, *domain.MetricsError) {

	db, err := sql.Open("pgx", config.DatabaseDsn)

	if err != nil {
		return nil, handleError(err.Error(), http.StatusInternalServerError)
	}

	storage := &Store{
		db: db,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {
		return nil, handleError(err.Error(), http.StatusInternalServerError)
	}

	if err := runMigrations(config); err != nil {
		return nil, err
	}

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
func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		logger.Log.Info("Error close db", zap.Error(err))
	}
}

func (s *Store) GetValues() ([]domain.Metrics, *domain.MetricsError) {
	metrics := make([]domain.Metrics, 0, 10)

	rows, err := s.db.QueryContext(context.Background(), "SELECT id, type, delta, value from metrics ")
	if err != nil {
		return nil, handleError(err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return nil, handleError(err.Error(), http.StatusInternalServerError)
		}
		metrics = append(metrics, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, handleError(err.Error(), http.StatusInternalServerError)
	}
	return metrics, nil
}
func (s *Store) GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError) {
	metrics := domain.Metrics{}

	row := s.db.QueryRowContext(context.Background(),
		"SELECT id, type, delta, value from metrics where id = $1 and type = $2 ", request.ID, request.MType)

	err := row.Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value)
	fmt.Println(err)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, handleError(err.Error(), http.StatusInternalServerError)
	}
	return &metrics, nil
}
func (s *Store) SetValue(request *domain.Metrics) *domain.MetricsError {
	_, err := s.db.Exec(
		"insert into metrics (id, type, delta, value) values ($1, $2, $3, $4)"+
			"on conflict (id) do update set id = excluded.id, type = excluded.type, delta = excluded.delta, value = excluded.value",
		request.ID, request.MType, request.Delta, request.Value,
	)

	if err != nil {
		return handleError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func runMigrations(config *Config) *domain.MetricsError {
	db, err := goose.OpenDBWithDriver("pgx", config.DatabaseDsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	if err := goose.RunContext(context.Background(), "up", db, config.MigrationsPath); err != nil {
		return handleError(err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func handleError(text string, code int) *domain.MetricsError {
	return &domain.MetricsError{
		Text: text,
		Code: code,
	}
}
