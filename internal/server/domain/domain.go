package domain

// MetricRequest is a request struct for metrics
type MetricRequest struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// UpdateResponse is a response struct for metrics update
type UpdateResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MetricResponse is a response struct for metrics
type MetricResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
