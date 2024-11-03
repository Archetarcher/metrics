package pgx

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
)

var (
	errConnectionException = errors.New("db connection exception")
	dbError                = 500
)

// Store is a struct for pgx storage, keeps configurations and sqlx.DB instance
type Store struct {
	db     *sqlx.DB
	config *Config
	stmts  struct {
		getMetric *sqlx.Stmt
		setMetric *sqlx.Stmt
	}
}

// NewStore creates pgx storage instance, runs migrations
func NewStore(config *Config, ctx context.Context) (*Store, *domain.MetricsError) {

	db := sqlx.MustOpen("pgx", config.DatabaseDsn)

	storage := &Store{
		db:     db,
		config: config,
	}

	err := storage.CheckConnection(ctx)

	if err != nil {
		return nil, err
	}

	if err := runMigrations(config, ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

// RetryConnection retries connection
func RetryConnection(error *domain.MetricsError, interval int, try int, config *Config, ctx context.Context) (*Store, *domain.MetricsError) {
	logger.Log.Info("retrying db connection", zap.Int("interval", interval), zap.Int("try", try))

	time.Sleep(time.Duration(interval) * time.Second)

	if try < 1 {
		logger.Log.Info("all attempts finished", zap.Int("interval", interval), zap.Int("try", try))
		return nil, error
	}

	var pgErr *pgconn.PgError

	if errors.As(error.Err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {

		s, err := NewStore(config, ctx)
		if err != nil {
			RetryConnection(err, interval+2, try-1, config, ctx)
		}
		if s != nil {
			logger.Log.Info("connection established", zap.Int("interval", interval), zap.Int("try", try))
			return s, nil
		}

	}
	return nil, error
}

// CheckConnection checks connection
func (s *Store) CheckConnection(ctx context.Context) *domain.MetricsError {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {

		return handleDBError(err, dbError)
	}
	return nil
}

// Close closes connection
func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		logger.Log.Info("Error close db", zap.Error(err))
	}
}

// GetValuesIn fetches metrics by keys in slice
func (s *Store) GetValuesIn(keys []string, ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	var metrics []domain.Metrics

	q, args, err := sqlx.In(metricsGetByKeyQuery, keys)
	if err != nil {
		return nil, handleDBError(err, dbError)
	}
	q = sqlx.Rebind(sqlx.DOLLAR, q)
	err = s.db.SelectContext(ctx, &metrics, q, args...)

	if err != nil {
		return nil, handleDBError(err, dbError)
	}

	return metrics, nil
}

// GetValues fetches all metrics
func (s *Store) GetValues(ctx context.Context) ([]domain.Metrics, *domain.MetricsError) {
	metrics := make([]domain.Metrics, 0, 10)

	rows, err := s.db.QueryContext(ctx, metricsGetAllQuery)
	if err != nil {
		return nil, handleDBError(err, dbError)
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return nil, handleDBError(err, dbError)
		}
		metrics = append(metrics, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, handleDBError(err, dbError)
	}
	return metrics, nil
}

// GetValue fetches metric data
func (s *Store) GetValue(request *domain.Metrics, ctx context.Context) (*domain.Metrics, *domain.MetricsError) {
	metrics := domain.Metrics{}

	row := s.db.QueryRowContext(ctx,
		metricsGetByIDAndTypeQuery, request.ID, request.MType)

	err := row.Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, handleDBError(err, dbError)
	}

	return &metrics, nil
}

// SetValue inserts or updates metric data
func (s *Store) SetValue(request *domain.Metrics, ctx context.Context) *domain.MetricsError {
	_, err := s.db.ExecContext(ctx,
		metricsCreateQuery,
		request.ID, request.MType, request.Delta, request.Value, getKey(*request),
	)
	if err != nil {
		return handleDBError(err, dbError)
	}

	return nil
}

// SetValues inserts or updates batch of metrics data
func (s *Store) SetValues(request []domain.Metrics, ctx context.Context) *domain.MetricsError {

	tx, err := s.db.Begin()
	if err != nil {
		return handleDBError(err, dbError)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		metricsCreateQuery)
	if err != nil {
		return handleDBError(err, dbError)
	}
	defer stmt.Close()

	for _, m := range request {
		_, err := stmt.ExecContext(ctx, m.ID, m.MType, m.Delta, m.Value, getKey(m))

		if err != nil {
			return handleDBError(err, dbError)
		}
	}
	err = tx.Commit()

	if err != nil {
		return handleDBError(err, dbError)
	}
	return nil
}
func getKey(request domain.Metrics) string {
	return request.ID + "_" + request.MType

}

func runMigrations(config *Config, ctx context.Context) *domain.MetricsError {
	db, err := goose.OpenDBWithDriver("pgx", config.DatabaseDsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	if err := goose.RunContext(ctx, "up", db, config.MigrationsPath); err != nil {
		return handleDBError(err, dbError)
	}

	return nil
}

func handleDBError(err error, code int) *domain.MetricsError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		err = errConnectionException
	}

	return &domain.MetricsError{
		Text: err.Error(),
		Code: code,
		Err:  err,
	}
}
