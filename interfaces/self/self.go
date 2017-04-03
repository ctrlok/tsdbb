package self

import "fmt"
import "github.com/ctrlok/tsdbb/interfaces"

type Metric struct{}

func (m *Metric) Name() string { return "" }

type TSDB struct {
	max    int
	metric Metric
}

func (t *TSDB) NewSender() interfaces.Sender {
	return &Sender{}
}
func (t *TSDB) GenerateMetrics(i int) { t.max = i }
func (t *TSDB) Metric(i int) (interfaces.Metric, error) {
	if i > t.max {
		return &t.metric, fmt.Errorf("error")
	}
	return &t.metric, nil
}

type Sender struct {
	host string
}

func (s *Sender) SetHost(host string) {
	s.host = host
}

func (s *Sender) GetHost() string {
	return s.host
}

func (s *Sender) Connect() error { return nil }

func (s *Sender) Send(metric interfaces.SendMetric) error {
	return nil
}
