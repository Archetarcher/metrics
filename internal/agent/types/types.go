package types

import "github.com/Archetarcher/metrics.git/internal/agent/domain"

// UpdateMetrics is a type for update metrics function.
type UpdateMetrics func(request []domain.Metrics) (*domain.SendResponse, *domain.MetricsError)
