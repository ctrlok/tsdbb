package server

import (
	"time"

	metrics "github.com/armon/go-metrics"
	i "github.com/ctrlok/tsdbb/interfaces"
)

type countStruct struct {
	count int
	step  int
}

func Loop(pregenerated i.PregeneratedMetrics, senders []i.Sender, count int, tickerChan <-chan time.Time, countChan chan countStruct) (err error) {
	metricsChan := make(chan i.SendMetric, 3000) // This is best value in my benchmarks
	for _, sender := range senders {
		go func(sender i.Sender) {
			for {
				senderInstance(sender, metricsChan)
			}
		}(sender)
	}
	err = tickerLoop(pregenerated, metricsChan, tickerChan, countChan, count)
	return
}

func LoopBench(pregenerated i.PregeneratedMetrics, senders []i.Sender, count int, tickerChan <-chan time.Time, countChan chan countStruct) (err error) {
	metricsChan := make(chan struct{}, 3000) // This is best value in my benchmarks
	for _, sender := range senders {
		go func(sender i.Sender) {
			n := 0
			for {
				<-metricsChan
				n++
				if n == 1000 {
					metrics.IncrCounter([]string{"s"}, 1)
					n = 0
				}

			}
		}(sender)
	}
	for {
		for i := 0; i < count; i++ {
			metricsChan <- struct{}{}
		}
	}
	return
}

type control struct {
	start int
	end   int
	time  *time.Time
}

func LoopBench2(pregenerated i.PregeneratedMetrics, senders []i.Sender, count int, tickerChan <-chan time.Time, countChan chan countStruct) (err error) {
	n := count / len(senders)
	countrolChan := make(chan control, 400)
	for _, sender := range senders {
		go func() {
			for {
				c := <-countrolChan
				for i := c.start; i < c.end; i++ {
					metric, _ := pregenerated.Metric(i)
					sender.Send(metric, c.time)
				}
				metrics.IncrCounter([]string{"s"}, float32(n))
			}
		}()
	}
	for t := range tickerChan {
		start := 0
		var end int
		for start < count {
			end = start + n
			countrolChan <- control{start: start, end: end, time: &t}
			start = end
		}
	}
	return
}

func LoopPool(pregenerated i.PregeneratedMetrics, senders chan i.Sender, count int, tickerChan <-chan time.Time, countChan chan countStruct) (err error) {
	newCount := countStruct{count: count, step: 0}
	for t := range tickerChan {
		select {
		case newCount = <-countChan:
			count = checkCount(count, &newCount)
			sendMetricsToPool(pregenerated, count, senders, &t)
		default:
			count = checkCount(count, &newCount)
			sendMetricsToPool(pregenerated, count, senders, &t)
		}
	}
	return
}

func sendMetricsToPool(pregenerated i.PregeneratedMetrics, count int, senders chan i.Sender, t *time.Time) {
	for n := 0; n < count; n++ {
		sender := <-senders
		go func() {
			metrics.IncrCounter([]string{"sender", "succes"}, 1)
			metric, _ := pregenerated.Metric(n)
			sender.Send(metric, t)
			senders <- sender
		}()
	}
}

func senderInstance(sender i.Sender, metricsChan chan i.SendMetric) {
	var err error
	metric, next := <-metricsChan
	if next {
		err = sender.Send(metric.Metric, metric.Time)
		if err != nil {
			metrics.IncrCounter([]string{"e"}, 1)
			return
		}
		metrics.IncrCounter([]string{"s"}, 1)
	}
}

func tickerLoop(pregenerated i.PregeneratedMetrics, metrics chan i.SendMetric, tickerChan <-chan time.Time, countChan chan countStruct, count int) (err error) {
	newCount := countStruct{count: count, step: 0}
	for t := range tickerChan {
		select {
		case newCount = <-countChan:
			count = checkCount(count, &newCount)
			err = sendMetricsToChannel(pregenerated, count, metrics, t)
			if err != nil {
				return
			}
		default:
			count = checkCount(count, &newCount)
			err = sendMetricsToChannel(pregenerated, count, metrics, t)
			if err != nil {
				return
			}
		}
	}
	return
}

func checkCount(initialCount int, newCount *countStruct) int {
	if newCount.count == initialCount {
		return initialCount
	} else if newCount.count > initialCount {
		tmpCount := initialCount + newCount.step
		if tmpCount > newCount.count {
			return newCount.count
		}
		return tmpCount
	}
	tmpCount := initialCount - newCount.step
	if tmpCount < newCount.count {
		return newCount.count
	}
	return tmpCount
}

func sendMetricsToChannel(pregenerated i.PregeneratedMetrics, count int, metrics chan i.SendMetric, t time.Time) (err error) {
	for n := 0; n < count; n++ {
		metric, err := pregenerated.Metric(n)
		if err != nil {
			return err
		}
		metrics <- i.SendMetric{Metric: metric, Time: &t}
	}
	return
}
