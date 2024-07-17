package domain

// Metrics struct
type Metrics struct {
	ID    string   `json:"id" db:"id"`                           // имя метрики
	MType string   `json:"type" db:"type"`                       // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value,omitempty"` // значение метрики в случае передачи gauge
}

// MetricsBatch struct
type MetricsBatch struct {
	Metrics []Metrics
}
