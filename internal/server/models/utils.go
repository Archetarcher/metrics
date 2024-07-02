package models

import "strconv"

const (
	EmptyParam  = ""
	GaugeType   = "gauge"
	CounterType = "counter"
)

var (
	RunAddr  string
	LogLevel string
)
var AllowedHeaders = map[string]string{
	//"Content-Type": "text/plain",
}

func GetStringValue(result *Metrics) string {
	value := ""
	if result.MType == GaugeType {
		value = strconv.FormatFloat(*result.Value, 'f', 3, 64)
	}
	if result.MType == CounterType {
		value = strconv.FormatInt(*result.Delta, 10)
	}

	return value
}
