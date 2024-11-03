package utils

import (
	"strconv"

	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

// GetStringValue converts metrics value to string
func GetStringValue(result *domain.Metrics) string {
	switch result.MType {
	case domain.GaugeType:
		return strconv.FormatFloat(*result.Value, 'f', 3, 64)
	case domain.CounterType:
		return strconv.FormatInt(*result.Delta, 10)
	default:
		return ""
	}
}
