package server

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	metrics "github.com/armon/go-metrics"
	i "github.com/ctrlok/tsdbb/interfaces"
)

var Logger, _ = zap.NewDevelopment()

func StartServer(pregenerated i.PregeneratedMetrics,
	senders []i.Sender, count int, tick, statTick time.Duration, listenURL string, statDisable bool) (err error) {
	var ticker = time.NewTicker(tick)
	var countChan = make(chan countStruct)

	var inm MetricSink
	if statDisable {
		inm = &BlackholeSink{}
	} else {
		inm = metrics.NewInmemSink(statTick, 3*statTick)
	}
	// TODO: send metrics to statsd, statsite, etc, change name
	metricsConfig := metrics.Config{
		ServiceName:          "s",
		EnableHostname:       false,
		EnableRuntimeMetrics: false,
		EnableTypePrefix:     false,
	}
	metrics.NewGlobal(&metricsConfig, inm)

	http.HandleFunc("/ShutDown", func(w http.ResponseWriter, r *http.Request) { shutDown(w, r, ticker) })
	go func() {
		err = http.ListenAndServe(listenURL, nil) // set listen port
		if err != nil {
			panic(err)
		}
	}()
	go logFunc(statTick, inm)

	err = loop(pregenerated, senders, count, ticker.C, countChan)
	if err != nil {
		Logger.Fatal(err.Error())
	}
	return

}

func logFunc(tick time.Duration, inm MetricSink) {
	ticker := time.NewTicker(tick)
	for range ticker.C {
		for _, metric := range inm.Data() {
			for k, v := range metric.Counters {
				Logger.Info(fmt.Sprintf("%s: %f", k, v.Sum))
			}
			// for k, v := range metric.Gauges {
			// 	logger.Info(fmt.Sprintf("%s: %v", k, v))
			// }
		}
		Logger.Info("-------------------------------")
	}
}

func shutDown(w http.ResponseWriter, r *http.Request, t *time.Ticker) {
	Logger.Info("Shutting down...")
	t.Stop()
}
