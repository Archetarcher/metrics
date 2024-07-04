package domain

const (
	EmptyParam  = ""
	GaugeType   = "gauge"
	CounterType = "counter"
)

var (
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	StoreInterval   int
	Restore         bool
)
