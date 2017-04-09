package server

import (
	"time"

	"strconv"

	metrics "github.com/armon/go-metrics"
	i "github.com/ctrlok/tsdbb/interfaces"
)

type countStruct struct {
	count int
	step  int
}

type control struct {
	start int
	end   int
	N     int
	time  []byte
}

var errorName = []string{"e"}
var succName = []string{"s"}

func loop(pregenerated i.PregeneratedMetrics, senders []i.Sender, count int, tickerChan <-chan time.Time, countChan chan countStruct) (err error) {
	controlChan := make(chan control, 400)
	for _, sender := range senders {
		go senderInstance(pregenerated, sender, controlChan)
	}
	tickerLoop(count, len(senders), tickerChan, controlChan, countChan)
	return
}

func senderInstance(pregenerated i.PregeneratedMetrics, sender i.Sender, controlChan chan control) (err error) {
	for c := range controlChan {
		for i := c.start; i < c.end; i++ {
			metric, err := pregenerated.Metric(i)
			if err != nil {
				return err
			}
			err = sender.Send(metric, c.time)
			if err != nil {
				metrics.IncrCounter(errorName, float32(c.N))
			}
			// metrics.IncrCounter([]string{"real_succ"}, 1)
		}
		metrics.IncrCounter(succName, float32(c.N))
	}
	return nil
}

func tickerLoop(count, senders int, tickerChan <-chan time.Time, controlChan chan control, countChan chan countStruct) {
	newCount := countStruct{count: count, step: 0}
	for t := range tickerChan {
		select {
		case newCount = <-countChan:
		default:
		}
		count = checkCount(count, &newCount)
		for _, c := range splitArray(count, senders, t) {
			controlChan <- c
		}
	}
}

// I know, it's really non weel performer. But that methot will work only 1 time/sec
// TODO: rewrite
func splitArray(count, senders int, t time.Time) (array []control) {
	if count <= 0 {
		return
	}
	n := count / senders
	timeByte := []byte(strconv.Itoa(int(t.Unix())))
	var i int
	for i = 0; i+n < count; i += n {
		array = append(array, control{start: i, end: i + n, N: n, time: timeByte})
	}
	if i == count {
		return
	}
	array = append(array, control{start: i, end: count, N: count - i, time: timeByte})
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
