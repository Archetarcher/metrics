package domain

// MetricsError is an error struct keeps code and message of error.
type MetricsError struct {
	Err  error
	Text string
	Code int
}
