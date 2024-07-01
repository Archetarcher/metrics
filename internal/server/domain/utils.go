package domain

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
