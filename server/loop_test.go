package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	i "github.com/ctrlok/tsdbb/interfaces"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testSender struct {
	sended int
	host   string
}

func (t *testSender) SetHost(string)  { return }
func (t *testSender) GetHost() string { return t.host }
func (t *testSender) Connect() error  { return nil }
func (t *testSender) Send(s i.SendMetric) error {
	t.sended++
	return nil
}

type testMetric struct{}

func (m *testMetric) Name() string {
	return ""
}

type testTSDB struct {
	getMetric int
}

func (t *testTSDB) NewSender() i.Sender   { return &testSender{} }
func (t *testTSDB) GenerateMetrics(i int) {}
func (t *testTSDB) Metric(i int) (m i.Metric, err error) {
	if i >= 10 {
		err = fmt.Errorf("bigger that exist")
		return &testMetric{}, err
	}
	t.getMetric++
	return &testMetric{}, nil
}

func TestSenderInstance(t *testing.T) {
	timeNow := time.Now()
	metric := testMetric{}
	sendMetric := i.SendMetric{Metric: &metric, Time: timeNow}

	// Do nothing if don't have any metric
	sender := testSender{sended: 0}
	metrics := make(chan i.SendMetric, 1)
	go senderInstance(&sender, metrics)
	time.Sleep(1 * time.Millisecond)
	assert.Zero(t, sender.sended, "Do nothing if don't have any metric")
	close(metrics)

	// Send message if it has messages in channel
	sender = testSender{sended: 0}
	metrics = make(chan i.SendMetric, 1)
	metrics <- sendMetric
	senderInstance(&sender, metrics)
	assert.NotZero(t, sender.sended, "Send message if it has messages in channel")
	assert.Equal(t, 1, sender.sended, "Send message if it has messages in channel")
	close(metrics)

	// Don't do anything if channel was closed
	sender = testSender{sended: 0}
	metrics = make(chan i.SendMetric, 1)
	close(metrics)
	senderInstance(&sender, metrics)
	assert.Zero(t, sender.sended, "Don't do anything if channel was closed")
}

func TestSendMetricsToChannel(t *testing.T) {
	var err error
	var tsdb *testTSDB
	var metrics chan i.SendMetric
	time := time.Now()

	metrics = make(chan i.SendMetric, 10)
	tsdb = &testTSDB{}
	err = sendMetricsToChannel(tsdb, 2, metrics, time)
	assert.NoError(t, err)
	assert.Equal(t, 2, tsdb.getMetric)
	assert.Equal(t, 2, len(metrics))
	close(metrics)

	metrics = make(chan i.SendMetric, 20)
	tsdb = &testTSDB{}
	err = sendMetricsToChannel(tsdb, 19, metrics, time)
	assert.Error(t, err)
	assert.Equal(t, 10, tsdb.getMetric)
	assert.Equal(t, 10, len(metrics))
	close(metrics)

}

func TestCheckCount(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(10, checkCount(10, &countStruct{count: 10, step: 2}))
	assert.Equal(10, checkCount(10, &countStruct{count: 10, step: 20}))
	assert.NotEqual(10, checkCount(10, &countStruct{count: 20, step: 20}))
	assert.Equal(12, checkCount(10, &countStruct{count: 12, step: 2}))
	assert.Equal(12, checkCount(10, &countStruct{count: 16, step: 2}))
	assert.Equal(12, checkCount(10, &countStruct{count: 12, step: 20}))
	assert.Equal(8, checkCount(10, &countStruct{count: 8, step: 2}))
	assert.Equal(8, checkCount(10, &countStruct{count: 6, step: 2}))
	assert.Equal(8, checkCount(10, &countStruct{count: 8, step: 20}))
}

func TestTickerLoop_Basic(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	close(tickerChan)

	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(metrics))
	assert.Equal(t, 2, tsdb.getMetric)
}

func TestTickerLoop_TwoTicks(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	close(tickerChan)

	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 4)
	assert.NoError(t, err)
	assert.Equal(t, 8, len(metrics))
	assert.Equal(t, 8, tsdb.getMetric)
}

func TestTickerLoop_ErrorOutOfIndex(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	close(tickerChan)

	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 12)
	assert.Error(t, err)
	assert.Equal(t, 10, len(metrics))
	assert.Equal(t, 10, tsdb.getMetric)
}

func TestTickerLoop_NoErrorOutOfIndex(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	close(tickerChan)
	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 6)
	assert.NoError(t, err)
	assert.Equal(t, 12, len(metrics))
	assert.Equal(t, 12, tsdb.getMetric)
}

func TestTickerLoop_countUp1(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	countChan <- countStruct{count: 10, step: 2}
	close(tickerChan)
	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 2)
	assert.NoError(t, err)
	// tick1: default 2 + 2 step = 4, tick2: 4 + 2 = 6, tick1 + tick2:
	assert.Equal(t, 10, len(metrics))
	assert.Equal(t, 10, tsdb.getMetric)
}

func TestTickerLoop_countUp2(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	countChan <- countStruct{count: 4, step: 2}
	close(tickerChan)
	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 2)
	assert.NoError(t, err)
	// tick1: default 2 + 2 step = 4, tick2: 4 + 0 = 4, tick1 + tick2:
	assert.Equal(t, 8, len(metrics))
	assert.Equal(t, 8, tsdb.getMetric)
}

func TestTickerLoop_countZero(t *testing.T) {
	var err error
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- timeNow
	tickerChan <- timeNow
	countChan <- countStruct{count: 0, step: 4}
	close(tickerChan)
	err = tickerLoop(&tsdb, metrics, tickerChan, countChan, 4)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(metrics))
	assert.Equal(t, 0, tsdb.getMetric)
}

func TestTickerLoop_countGoroutine(t *testing.T) {
	done := make(chan error)
	timeNow := time.Now()
	tsdb := testTSDB{}
	metrics := make(chan i.SendMetric, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	go func() {
		done <- tickerLoop(&tsdb, metrics, tickerChan, countChan, 4)
	}()

	tickerChan <- timeNow
	for len(tickerChan) != 0 {
	} // Wait until goroutine done
	for len(metrics) == 0 {
	}
	time.Sleep(1 * time.Millisecond)
	assert.Zero(t, len(done))
	assert.Equal(t, 4, len(metrics))
	assert.Equal(t, 4, tsdb.getMetric)

	for len(metrics) != 0 {
		<-metrics
	}
	tsdb.getMetric = 0

	countChan <- countStruct{count: 8, step: 2}
	tickerChan <- timeNow
	for len(tickerChan) != 0 {
	} // Wait until goroutine done
	for len(metrics) == 0 {
	}
	time.Sleep(1 * time.Millisecond)
	assert.Zero(t, len(done))
	assert.Equal(t, 6, len(metrics))
	assert.Equal(t, 6, tsdb.getMetric)

	for len(metrics) != 0 {
		<-metrics
	}
	tsdb.getMetric = 0

	tickerChan <- timeNow
	for len(tickerChan) != 0 {
	} // Wait until goroutine done
	for len(metrics) == 0 {
	}
	time.Sleep(1 * time.Millisecond)
	assert.Zero(t, len(done))
	assert.Equal(t, 8, len(metrics))
	assert.Equal(t, 8, tsdb.getMetric)

	for len(metrics) != 0 {
		<-metrics
	}
	tsdb.getMetric = 0

	tickerChan <- timeNow
	for len(tickerChan) != 0 {
	} // Wait until goroutine done
	for len(metrics) == 0 {
	}
	time.Sleep(1 * time.Millisecond)
	assert.Zero(t, len(done))
	assert.Equal(t, 8, len(metrics))
	assert.Equal(t, 8, tsdb.getMetric)

	close(tickerChan)
}

func BenchmarkCheckCount_eq(b *testing.B) {
	newCount := countStruct{count: 10, step: 20}
	for n := 0; n < b.N; n++ {
		_ = checkCount(10, &newCount)
	}
}

func BenchmarkCheckCount_hi(b *testing.B) {
	newCount := countStruct{count: 100, step: 20}
	for n := 0; n < b.N; n++ {
		_ = checkCount(10, &newCount)
	}
}

type benchSender struct {
	host string
}

func (t *benchSender) SetHost(string)            { return }
func (t *benchSender) GetHost() string           { return t.host }
func (t *benchSender) Connect() error            { return nil }
func (t *benchSender) Send(s i.SendMetric) error { return nil }

type benchMetric struct{}

func (m *benchMetric) Name() string { return "" }

type benchTSDB struct {
	metric i.Metric
}

func (t *benchTSDB) NewSender() i.Sender   { return &benchSender{} }
func (t *benchTSDB) GenerateMetrics(i int) {}
func (t *benchTSDB) Metric(i int) (m i.Metric, err error) {
	return t.metric, nil
}

func BenchmarkSendMetricsToChannel(b *testing.B) {
	time := time.Now()
	metrics := make(chan i.SendMetric, 100)
	tsdb := &benchTSDB{metric: &benchMetric{}}

	go func() {
		for {
			<-metrics
		}
	}()

	for n := 0; n < b.N; n++ {
		sendMetricsToChannel(tsdb, 1, metrics, time)
	}
}

// Send metrics to blackhole is much faster than incerement metrics.
// Anyway, it was really faster without metrics, but we need statistics.
//
// BenchmarkSenderInstance-4   										50000000	       350 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSenderInstance_MetricsBlackhole-4   	30000000	       526 ns/op	      48 B/op	       1 allocs/op
// BenchmarkSenderInstance_MetricsInmem-4       	 5000000	      2593 ns/op	     208 B/op	       4 allocs/op
func BenchmarkSenderInstance_MetricsBlackhole(b *testing.B) {
	sender := &benchSender{host: "host"}
	metricChan := make(chan i.SendMetric, 10000)
	timeNow := time.Now()
	go func() {
		for {
			metricChan <- i.SendMetric{Metric: &benchMetric{}, Time: timeNow}
		}
	}()

	for n := 0; n < b.N; n++ {
		senderInstance(sender, metricChan)
	}
}

func BenchmarkSenderInstance_MetricsInmem(b *testing.B) {
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.NewGlobal(metrics.DefaultConfig("service-name"), inm)
	sender := &benchSender{host: "host"}
	metricChan := make(chan i.SendMetric, 10000)
	timeNow := time.Now()
	go func() {
		for {
			metricChan <- i.SendMetric{Metric: &benchMetric{}, Time: timeNow}
		}
	}()

	for n := 0; n < b.N; n++ {
		senderInstance(sender, metricChan)
	}
	metrics.NewGlobal(metrics.DefaultConfig("service-name"), &metrics.BlackholeSink{})
}

func BenchmarkTickerLoop(b *testing.B) {
	timeNow := time.Now()
	tsdb := &benchTSDB{metric: &benchMetric{}}
	metrics := make(chan i.SendMetric, 100)
	countChan := make(chan countStruct, 100)
	tickerChan := make(chan time.Time, 1)

	go func() {
		for {
			<-metrics
		}
	}()

	go tickerLoop(tsdb, metrics, tickerChan, countChan, 1)
	for n := 0; n < b.N; n++ {
		tickerChan <- timeNow
	}
}

func BenchmarkLoop(b *testing.B) {
	timeNow := time.Now()
	tsdb := &benchTSDB{metric: &benchMetric{}}
	tickerChan := make(chan time.Time)
	senders := []i.Sender{
		&benchSender{},
		&benchSender{},
		&benchSender{},
		&benchSender{},
	}

	go Loop(tsdb, senders, 1000, tickerChan)
	for n := 0; n < b.N; n++ {
		tickerChan <- timeNow
	}

}

func BenchmarkCheckCount_lo(b *testing.B) {
	newCount := countStruct{count: 10, step: 20}
	for n := 0; n < b.N; n++ {
		_ = checkCount(100, &newCount)
	}
}

func BenchmarkLogger_NoLogger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = n
	}
}

func BenchmarkLogger_LoggerNotShow(b *testing.B) {
	logger, _ := zap.NewProduction()
	for n := 0; n < b.N; n++ {
		logger.Debug("")
		logger.Debug("")
	}
}

func BenchmarkMetrics_Blackhole(b *testing.B) {
	for n := 0; n < b.N; n++ {
		metrics.IncrCounter([]string{"key"}, 1)
	}
}

func BenchmarkMetrics_Inmem(b *testing.B) {
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	for n := 0; n < b.N; n++ {
		inm.IncrCounter([]string{"key"}, 1)
		inm.IncrCounter([]string{"key2"}, 1)
		inm.IncrCounter([]string{"key3"}, 1)
	}
}
