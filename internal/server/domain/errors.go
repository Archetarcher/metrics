package domain

// MetricsError is an error struct keeps code, message of error and general error interface{}.
type MetricsError struct {
	Err  error
	Text string
	Code int
}
