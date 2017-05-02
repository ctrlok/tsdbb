package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/ctrlok/tsdbb/log"
)

func StartServer(basic interfaces.Basic, opts Options, ctx context.Context) {
	timeNow := time.Now().UnixNano()
	basic.NewRequests(100000000)
	log.SLogger.Infow("Metrics generated", "timer_ns", int((time.Now().UnixNano()-timeNow)/1000000))
	bus := make(chan busMessage, 10000000)
	chStat := make(chan statMessage)
	go statisctics(ctx, chStat)

	err := startClients(ctx, basic, opts, bus, chStat)
	if err != nil {
		log.Logger.Error("Fail to start server!", log.ParseFields(ctx)...)
	}

	controlChan := make(chan ControlMessages, 1)
	tickChan := time.NewTicker(opts.Tick)

	go startGenerator(ctx, opts, controlChan, bus, tickChan.C)

	log.SLogger.Info("Starting server")
	log.Logger.Info("Starting server at port: ", log.ParseFields(ctx)...)
	http.HandleFunc("/upto", upto(controlChan))
	err = http.ListenAndServe(opts.ListenURL, nil) // set listen port
	if err != nil {
		panic(err)
	}
}

func upto(ch chan ControlMessages) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg ControlMessages
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			log.Logger.Error("Error setting new count and step")
			w.WriteHeader(http.StatusInternalServerError)
		}
		ch <- msg
		w.WriteHeader(http.StatusOK)
	}
}
