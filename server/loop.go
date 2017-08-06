package server

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/ctrlok/tsdbb/log"
	"github.com/rcrowley/go-metrics"
)

type busMessage struct {
	ctx      context.Context
	start, N int
	time     []byte
}

type ControlMessages struct {
	Count, Step int
}

type Options struct {
	StartCount, StartStep int
	Parallel              int
	Tick, StatTick        time.Duration
	ListenURL             string
	Servers               []string
	StatPrefix            string
}

func startClients(ctx context.Context, basic interfaces.Basic, opts Options,
	bus <-chan busMessage) (err error) {
	ctx = context.WithValue(ctx, log.KeyOperation, "startClients")
	if len(opts.Servers) == 0 {
		err = fmt.Errorf("You should define at least one server to send metrics")
		log.Logger.Error(err.Error(), log.ParseFields(ctx)...)
		return err
	}
	for _, server := range opts.Servers {
		for i := 0; i < opts.Parallel; i++ {
			uri, err := url.Parse(server)
			if err != nil {
				log.Logger.Error("Error parsing url "+server, log.ParseFields(ctx)...)
				return err
			}
			ctx2 := context.WithValue(ctx, log.KeyUrl, uri.Path)
			ctx2 = context.WithValue(ctx2, log.KeyClientNum, i)
			cli, err := basic.NewClient(uri)
			if err != nil {
				log.Logger.Error("Error creating client "+err.Error(), log.ParseFields(ctx2)...)
				return err
			}
			go startClient(ctx2, cli, basic, bus)
		}
	}
	return nil
}

func startClient(ctx context.Context, cli interfaces.Client, basic interfaces.Basic,
	bus <-chan busMessage) {
	ctx = context.WithValue(ctx, log.KeyOperation, "sendMessage")
	log.Logger.Debug("Starting...", log.ParseFields(ctx)...)
	metricNameHostSucc := metrics.NewCounter()
	metrics.Register(fmt.Sprintf("benchcli.%s.sended", cli.Host()), metricNameHostSucc)
	metricNameAllSucc := metrics.NewCounter()
	metrics.Register("benchcli.all.sended", metricNameAllSucc)
	metricNameHostErr := metrics.NewCounter()
	metrics.Register(fmt.Sprintf("benchcli.%s.error", cli.Host()), metricNameHostErr)
	metricNameAllErr := metrics.NewCounter()
	metrics.Register("benchcli.all.error", metricNameAllErr)
	for message := range bus {
		end := message.start + message.N
	SEND_MESSAGES:
		for n := message.start; n < end; n++ {
			select {
			case <-message.ctx.Done():
				break SEND_MESSAGES
			default:
				metricNameHostSucc.Inc(1)
				metricNameAllSucc.Inc(1)
				err := cli.Send(basic.Req(n), message.time)
				if err != nil {
					log.Logger.Debug("Error sending message", log.ParseFields(ctx)...)
					metricNameHostErr.Inc(1)
					metricNameAllErr.Inc(1)
					break
				}
			}
		}
	}
}

func startGenerator(ctx context.Context, opts Options,
	controlChan <-chan ControlMessages, bus chan<- busMessage, tickChan <-chan time.Time) {
	ctx = context.WithValue(ctx, log.KeyOperation, "startGenerator")
	defaultControl := ControlMessages{opts.StartCount, opts.StartStep}
	count := opts.StartCount
	for t := range tickChan {
		select {
		case tmpControl := <-controlChan:
			if tmpControl.Step == 0 {
				tmpControl.Step = defaultControl.Step
			}
			defaultControl = tmpControl
		default:
		}
		count = checkCount(count, &defaultControl)
		ctx2 := context.WithValue(ctx, log.KeyPlannedCount, count)
		ctx3, _ := context.WithTimeout(ctx2, opts.Tick)
		for _, c := range splitArray(ctx3, count, opts.Parallel*len(opts.Servers), t) {
			bus <- c
		}
	}
}

func checkCount(initialCount int, defaultControl *ControlMessages) int {
	switch {
	case defaultControl.Count == initialCount:
		return initialCount
	case defaultControl.Count > initialCount:
		tmpCount := initialCount + defaultControl.Step
		if tmpCount < defaultControl.Count {
			return tmpCount
		}
		return defaultControl.Count
	}
	tmpCount := initialCount - defaultControl.Step
	if tmpCount > defaultControl.Count {
		return tmpCount
	}
	return defaultControl.Count
}

// I know, it's really non well performer. But that method will work only 1 time/sec
// TODO: rewrite
func splitArray(ctx context.Context, count, senders int, t time.Time) (array []busMessage) {
	if count <= 0 {
		return
	}
	n := count / senders
	timeByte := []byte(fmt.Sprint(strconv.Itoa(int(t.Unix()))))
	var k int
	for k = 0; k+n < count; k += n {
		array = append(array, busMessage{start: k, N: n, time: timeByte, ctx: ctx})
	}
	array = append(array, busMessage{start: k, N: count - k, time: timeByte, ctx: ctx})
	return
}
