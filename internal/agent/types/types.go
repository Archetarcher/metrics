package types

import "github.com/Archetarcher/metrics.git/internal/agent/domain"

type UpdateMetrics func(request []domain.Metrics) (*domain.SendResponse, *domain.MetricsError)
