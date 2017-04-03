package main

import "time"
import "github.com/armon/go-metrics"

func loop(tsdb TSDB, senders []Sender, count int, tickerChan chan time.Time) (err error) {
	metricsChan := make(chan SendMetric, 100000) // This is best value in my benchmarks
	// metrics := make(chan SendMetric, len(senders))
	defer close(metricsChan)
	for _, sender := range senders {
		go func(sender Sender) {
			for {
				senderInstance(sender, metricsChan)
			}
		}(sender)
	}
	for t := range tickerChan {
		err = sendMetricsToChannel(tsdb, count, metricsChan, t)
		if err != nil {
			return err
		}
	}
	return
}

func senderInstance(sender Sender, metricsChan chan SendMetric) {
	metric, next := <-metricsChan
	if next {
		err := sender.Send(metric)
		if err != nil {
			metrics.IncrCounter([]string{"sender", sender.GetHost(), "error"}, 1)
		}
		metrics.IncrCounter([]string{"sender", sender.GetHost(), "succes"}, 1)
	} else {
		return
	}
}

type countStruct struct {
	count int
	step  int
}

func tickerLoop(tsdb TSDB, metrics chan SendMetric, tickerChan chan time.Time, countChan chan countStruct, count int) (err error) {
	newCount := countStruct{count: count, step: 0}
	for t := range tickerChan {
		select {
		case newCount = <-countChan:
			count = checkCount(count, &newCount)
			err = sendMetricsToChannel(tsdb, count, metrics, t)
			if err != nil {
				return
			}
		default:
			count = checkCount(count, &newCount)
			err = sendMetricsToChannel(tsdb, count, metrics, t)
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

func sendMetricsToChannel(tsdb TSDB, count int, metrics chan SendMetric, t time.Time) (err error) {
	for i := 0; i < count; i++ {
		metric, err := tsdb.Metric(i)
		if err != nil {
			return err
		}
		metrics <- SendMetric{Metric: metric, Time: t}
	}
	return
}
