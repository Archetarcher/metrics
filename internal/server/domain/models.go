package domain

// Metrics is a metric struct keeps type, name, and value of metrics.
type Metrics struct {
	ID    string   `json:"id" db:"id"`                           // metrics name
	MType string   `json:"type" db:"type"`                       // metrics type, accepts value gauge or counter.
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty"` // metrics value if provided type is counter.
	Value *float64 `json:"value,omitempty" db:"value,omitempty"` // metrics value if provided type is gauge.
}
