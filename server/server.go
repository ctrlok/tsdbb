package server

import (
	"context"

	"net/http"

	"time"

	"github.com/ctrlok/tsdbb-2/interfaces"
	"github.com/ctrlok/tsdbb-2/log"
)

func StartServer(basic interfaces.Basic, opts Options, ctx context.Context) {
	timeNow := time.Now().UnixNano()
	basic.NewRequests(100000000)
	log.SLogger.Infow("Metrics generated", "timer_ns", int((time.Now().UnixNano()-timeNow)/1000000))
	bus := make(chan busMessage, 10000000)
	chStat := make(chan statMessage)
	err := startClients(ctx, basic, opts, bus, chStat)
	go statisctics(ctx, chStat)
	if err != nil {
		log.Logger.Error("Fail to start server!", log.ParseFields(ctx)...)
	}
	controlChan := make(chan controlMessages, 1)
	go startGenerator(ctx, opts, controlChan, bus)

	log.SLogger.Info("Starting server")
	log.Logger.Info("Starting server at port: ", log.ParseFields(ctx)...)
	err = http.ListenAndServe(opts.ListenURL, nil) // set listen port
	if err != nil {
		panic(err)
	}
}
