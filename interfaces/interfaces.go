package interfaces

import (
	"time"
)

// SendMetric is a struct which contain Metric and time of current tick.
// Use SendMetric.Time in Sender.Send setting time.
type SendMetric struct {
	Metric Metric
	Time   time.Time
}

// Metric is a main interface for metrics
type Metric interface {
	Name() string
}

// TSDB can Generate metrics into internal format and show metric
type TSDB interface {
	GenerateMetrics(int)
	Metric(i int) (Metric, error)
	NewSender() Sender
}

// Sender is a main interface for sending metrics
type Sender interface {
	SetHost(string)
	GetHost() string
	Connect() error
	Send(SendMetric) error
}
