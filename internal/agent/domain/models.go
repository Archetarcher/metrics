package domain

// SendResponse is a response struct of server response after sending metrics.
type SendResponse struct {
	Status int `json:"status"`
}

// Metrics is a metric struct keeps type, name, and value of metrics.
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // metrics value if provided type is counter.
	Value *float64 `json:"value,omitempty"` // metrics value if provided type is gauge.
	ID    string   `json:"id"`              // metrics name
	MType string   `json:"type"`            // metrics type, accepts value gauge or counter.

}

// Gauge is a struct with all fields of metrics with type gauge.
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
	RandomValue,
	TotalMemory,
	FreeMemory,
	CPUutilization1 float64
}

// MetricsData is a map of Metrics.
type MetricsData map[string]Metrics
