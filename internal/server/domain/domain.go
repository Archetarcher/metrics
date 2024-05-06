package domain

// UpdateRequest is a request struct for services.service.update
type UpdateRequest struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// UpdateResponse is a response struct for services.service.update
type UpdateResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MetricResponse is a response struct
type MetricResponse struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
