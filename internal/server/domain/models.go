package domain

// Metrics is a metric struct keeps type, name, and value of metrics.
type Metrics struct {
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty"` // metrics value if provided type is counter.
	Value *float64 `json:"value,omitempty" db:"value,omitempty"` // metrics value if provided type is gauge.
	ID    string   `json:"id" db:"id"`                           // metrics name
	MType string   `json:"type" db:"type"`                       // metrics type, accepts value gauge or counter.
}

// SessionRequest is a session struct keeps encrypted key
type SessionRequest struct {
	Key []byte `json:"key"`
}
