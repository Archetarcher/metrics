package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

type Store interface {
	GetValues() ([]domain.Metrics, error)
	GetValue(request *domain.Metrics) (*domain.Metrics, error)
	SetValue(request *domain.Metrics) error
	CheckConnection() *domain.MetricsError
}

func NewStore(conf Config) (Store, *domain.MetricsError) {
	if conf.Memory.Active {
		return memory.NewStore(conf.Memory), nil
	}
	if conf.Pgx.Active {
		return pgx.NewStore(conf.Pgx)
	}

	return nil, nil
}
