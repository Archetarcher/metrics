package domain

// MetricData is a request struct for services.service.update
type MetricData struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// ServerResponse is a response struct
type ServerResponse struct {
	Status int `json:"status"`
}
