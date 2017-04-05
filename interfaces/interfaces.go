package interfaces

import (
	"net/url"
	"time"
)

// SendMetric is a struct which contain Metric and time of current tick.
// Use SendMetric.Time in Sender.Send setting time.
type SendMetric struct {
	Metric Metric
	Time   *time.Time
}

// Metric is a main interface for metrics
type Metric interface {
	Name() interface{}
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
	Send(Metric, *time.Time) error
}
