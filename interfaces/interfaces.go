package interfaces

import "net/url"

// Metric is a main interface for metrics
type Metric interface {
	Internal() interface{}
}

// PregeneratedMetrics can Generate metrics into internal format and show metric
type PregeneratedMetrics interface {
	Metric(i int) (Metric, error)
}

type TSDB interface {
	GenerateMetrics(int) PregeneratedMetrics
	NewSender(*url.URL) (Sender, error)
}

// Sender is a main interface for sending metrics
type Sender interface {
	GetHost() string
	Send(Metric, []byte) error
}
