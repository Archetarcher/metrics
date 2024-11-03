package domain

// MetricsError is an error struct keeps code, message of error and general error interface{}.
type MetricsError struct {
	Text string
	Code int
	Err  error
}
