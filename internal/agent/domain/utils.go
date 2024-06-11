package domain

type Gauge struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	TotalAlloc,
	RandomValue float64
}
type Counter struct {
	PollCount int64
}

var CounterMetric = "PollCount"

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

var (
	FlagServerRunAddr string
	ReportInterval    int
	PollInterval      int
)
