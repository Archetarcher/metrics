package domain

// MetricsError is an error struct keeps code and message of error.
type MetricsError struct {
	Text string
	Code int
	Err  error
}
