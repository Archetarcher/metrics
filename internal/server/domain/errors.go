package domain

type MetricsError struct {
	Text string
	Code int
	Err  error
}
