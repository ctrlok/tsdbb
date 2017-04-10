package server

import (
	"fmt"
	"strconv"
	"time"

	metrics "github.com/armon/go-metrics"
	i "github.com/ctrlok/tsdbb/interfaces"
	"go.uber.org/zap"
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
		for n := c.start; n < c.end; n++ {
			metric, err := pregenerated.Metric(n)
			if err != nil {
				return err
			}
			err = sender.Send(metric, c.time)
			if err != nil {
				Logger.Error("Finish sender")
				metrics.IncrCounter(errorName, float32(c.N))
			}
			// metrics.IncrCounter([]string{"real_succ"}, 1)
		}
		Logger.Debug("messages sended", zap.Int("sended", c.N), zap.String("to_host", sender.GetHost()))
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
	timeByte := []byte(fmt.Sprint(strconv.Itoa(int(t.Unix())), "\n"))
	var k int
	for k = 0; k+n < count; k += n {
		array = append(array, control{start: k, end: k + n, N: n, time: timeByte})
	}
	if k == count {
		return
	}
	array = append(array, control{start: k, end: count, N: count - k, time: timeByte})
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
