package server

import "time"
import "github.com/armon/go-metrics"
import i "github.com/ctrlok/tsdbb/interfaces"

func loop(tsdb i.TSDB, senders []i.Sender, count int, tickerChan chan time.Time) (err error) {
	metricsChan := make(chan i.SendMetric, 100000) // This is best value in my benchmarks
	// metrics := make(chan SendMetric, len(senders))
	defer close(metricsChan)
	for _, sender := range senders {
		go func(sender i.Sender) {
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

func senderInstance(sender i.Sender, metricsChan chan i.SendMetric) {
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

func tickerLoop(tsdb i.TSDB, metrics chan i.SendMetric, tickerChan chan time.Time, countChan chan countStruct, count int) (err error) {
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

func sendMetricsToChannel(tsdb i.TSDB, count int, metrics chan i.SendMetric, t time.Time) (err error) {
	for n := 0; n < count; n++ {
		metric, err := tsdb.Metric(n)
		if err != nil {
			return err
		}
		metrics <- i.SendMetric{Metric: metric, Time: t}
	}
	return
}
