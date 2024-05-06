package domain

const (
	EmptyParam  = ""
	GaugeType   = "gauge"
	CounterType = "counter"
)

var AllowedHeaders = map[string]string{
	//"Content-Type": "text/plain",
}

type Gauge struct {
	Name  string
	Value float64
}

type Counter struct {
	Name  string
	Value int64
}
