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

func (t *testSender) GetHost() string { return t.host }
func (t *testSender) Send(metric i.Metric, time *time.Time) error {
	t.sended++
	return nil
}

type testMetric struct{}

func (m *testMetric) Name() interface{} {
	return ""
}

type testPregeneratedMetrics struct {
	getMetric int
}

func (t *testPregeneratedMetrics) Metric(i int) (m i.Metric, err error) {
	if i >= 10 {
		err = fmt.Errorf("bigger that exist")
		return &testMetric{}, err
	}
	t.getMetric++
	return &testMetric{}, nil
}

func TestSenderInstance_Empty(t *testing.T) {
	pregenerated := testPregeneratedMetrics{}
	controlChan := make(chan control, 1)
	sender := testSender{}

	// Do nothing if don't have any metric
	close(controlChan)
	err := senderInstance(&pregenerated, &sender, controlChan)
	assert.NoError(t, err)
	assert.Equal(t, 0, sender.sended)
	assert.Equal(t, 0, pregenerated.getMetric)
}

func TestSenderInstance_Succ(t *testing.T) {
	time := time.Now()
	pregenerated := testPregeneratedMetrics{}
	controlChan := make(chan control, 2)
	sender := testSender{}

	// Send message if it has messages in channel
	controlChan <- control{start: 0, end: 2, N: 2, time: &time}
	controlChan <- control{start: 2, end: 4, N: 2, time: &time}
	close(controlChan)
	err := senderInstance(&pregenerated, &sender, controlChan)
	assert.NoError(t, err)
	assert.Equal(t, 4, sender.sended)
	assert.Equal(t, 4, pregenerated.getMetric)
}

func TestSenderInstance_Fail(t *testing.T) {
	time := time.Now()
	pregenerated := testPregeneratedMetrics{}
	controlChan := make(chan control, 2)
	sender := testSender{}

	// Send message if it has messages in channel
	controlChan <- control{start: 0, end: 2, N: 2, time: &time}
	controlChan <- control{start: 2, end: 50, N: 2, time: &time}
	close(controlChan)
	err := senderInstance(&pregenerated, &sender, controlChan)
	assert.Error(t, err)
	assert.Equal(t, 10, sender.sended)
	assert.Equal(t, 10, pregenerated.getMetric)
}

func TestSplitArray_Long(t *testing.T) {
	for count := 100; count < 400; count++ {
		for senders := 2; senders <= count; senders++ {
			array := splitArray(count, senders, time.Now())
			testArray := []int{}
			for _, c := range array {
				for start := c.start; start < c.end; start++ {
					testArray = append(testArray, start)
				}
			}
			if len(testArray) != count {
				t.Fatalf("non equal len. count: %v, senders: %v\n   testArray: %#v\n   array: %#v", count, senders, testArray, array)
			}

			for i := 0; i < count; i++ {
				if testArray[i] != i {
					t.Fatalf("Error splitting array, when count=%v, senders=%v, on %v element (actual: %v)", count, senders, i, testArray[i])
				}
			}
		}
	}
}

func TestSplitArray_N(t *testing.T) {
	array := splitArray(100, 3, time.Now())
	assert.Equal(t, 4, len(array))
	assert.Equal(t, 1, array[3].N)

	array = splitArray(99, 3, time.Now())
	assert.Equal(t, 3, len(array))
	assert.Equal(t, 33, array[2].N)
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

func TestTickerLoop_Skip(t *testing.T) {
	control := make(chan control, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	close(tickerChan)
	tickerLoop(100, 3, tickerChan, control, countChan)
	assert.Zero(t, len(control))
}

func TestTickerLoop_Basic(t *testing.T) {
	control := make(chan control, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- time.Now()
	close(tickerChan)
	tickerLoop(100, 3, tickerChan, control, countChan)
	assert.NotZero(t, len(control))
	assert.Equal(t, 4, len(control), "There should be 4 message with 33,33,33,1 elements to make")
}

func TestTickerLoop_TwoTicks(t *testing.T) {
	control := make(chan control, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- time.Now()
	tickerChan <- time.Now()
	close(tickerChan)
	tickerLoop(100, 3, tickerChan, control, countChan)
	assert.NotZero(t, len(control))
	assert.Equal(t, 8, len(control), "There should be 4 message with 33,33,33,1 elements to make")
}

func TestTickerLoop_CountUp1(t *testing.T) {
	control := make(chan control, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- time.Now()
	tickerChan <- time.Now()
	countChan <- countStruct{count: 10, step: 2}
	close(tickerChan)
	tickerLoop(1, 1, tickerChan, control, countChan)
	assert.NotZero(t, len(control))
	assert.Equal(t, 2, len(control), "")
	c1 := <-control
	assert.Equal(t, 3, c1.N)
	c2 := <-control
	assert.Equal(t, 5, c2.N)
}

func TestTickerLoop_Zero(t *testing.T) {
	control := make(chan control, 100)
	tickerChan := make(chan time.Time, 100)
	countChan := make(chan countStruct, 100)

	tickerChan <- time.Now()
	tickerChan <- time.Now()
	countChan <- countStruct{count: 0, step: 2}
	close(tickerChan)
	tickerLoop(1, 1, tickerChan, control, countChan)
	assert.Zero(t, len(control))
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

func (t *benchSender) GetHost() string           { return t.host }
func (t *benchSender) Send(s i.SendMetric) error { return nil }

type benchMetric struct{}

func (m *benchMetric) Name() interface{} { return "" }

type benchPregeneratedMetrics struct {
	metric i.Metric
}

func (t *benchPregeneratedMetrics) Metric(i int) (m i.Metric, err error) {
	return t.metric, nil
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
