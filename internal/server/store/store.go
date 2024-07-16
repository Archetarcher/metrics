package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

type Store interface {
	GetValues() ([]domain.Metrics, *domain.MetricsError)
	GetValue(request *domain.Metrics) (*domain.Metrics, *domain.MetricsError)
	SetValue(request *domain.Metrics) *domain.MetricsError
	CheckConnection() *domain.MetricsError
	Close()
}

func NewStore(conf Config) (Store, *domain.MetricsError) {

	if conf.Pgx.DatabaseDsn != domain.EmptyParam {
		return pgx.NewStore(conf.Pgx)
	}

	return memory.NewStore(conf.Memory), nil

}
