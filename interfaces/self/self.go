package self

import (
	"fmt"
	"net/url"

	"github.com/ctrlok/tsdbb/interfaces"
)

type Metric struct{}

func (m *Metric) Name() string { return "" }

type PregeneratedMetrics struct {
	max    int
	metric Metric
}

func (t *PregeneratedMetrics) Metric(i int) (interfaces.Metric, error) {
	if i > t.max {
		return &t.metric, fmt.Errorf("error")
	}
	return &t.metric, nil
}

type TSDB struct{}

func (t *TSDB) NewSender(uri *url.URL) (interfaces.Sender, error) {
	sender := &Sender{}
	sender.host = uri.Host
	return sender, nil
}
func (t *TSDB) GenerateMetrics(i int) interfaces.PregeneratedMetrics {
	return &PregeneratedMetrics{max: i}
}

type Sender struct {
	host string
}

func (s *Sender) GetHost() string {
	return s.host
}

func (s *Sender) Send(metric interfaces.SendMetric) error {
	return nil
}
