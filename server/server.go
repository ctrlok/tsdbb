package server

import (
	"time"

	i "github.com/ctrlok/tsdbb/interfaces"
)

func StartServer(pregenerated i.PregeneratedMetrics, senders []i.Sender,
	count int, tick time.Duration) (err error) {

	ticker := time.NewTicker(tick)
	countChan := make(chan countStruct)
	err = Loop(pregenerated, senders, count, ticker.C, countChan)
	return

}
