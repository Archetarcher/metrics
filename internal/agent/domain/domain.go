package domain

// MetricData is a data struct for metrics
type MetricData struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// SendResponse is a response struct for send response
type SendResponse struct {
	Status int `json:"status"`
}
