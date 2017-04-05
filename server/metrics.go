package server

import (
	"github.com/armon/go-metrics"
)

type MetricSink interface {
	metrics.MetricSink
	Data() []*metrics.IntervalMetrics
}

type BlackholeSink struct {
	metrics.BlackholeSink
}

func (*BlackholeSink) Data() (i []*metrics.IntervalMetrics) {
	return i
}
