package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/ctrlok/tsdbb/log"
	"github.com/ctrlok/tsdbb/server"
	"github.com/spf13/cobra"
)

var startCount int
var parallel int
var tick time.Duration
var maxMetrics int
var statDisable bool
var statTick time.Duration

// BenchCmd represents the bench command
var BenchCmd = &cobra.Command{
	Use:   "bench",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	RootCmd.AddCommand(BenchCmd)
	BenchCmd.PersistentFlags().IntVarP(&startCount, "start-count", "c", 1000, "Start count of metrics which will be send")
	BenchCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 6, "count of workers for each server in args")
	BenchCmd.PersistentFlags().DurationVarP(&tick, "tick", "t", 1*time.Second, "retention period")
	BenchCmd.PersistentFlags().IntVar(&maxMetrics, "maximum-metrics", 100000000, "maximum of metrics, which would be sended for one tick")
	BenchCmd.PersistentFlags().BoolVar(&statDisable, "statistics-disable", false, "disable internal metrics")
	BenchCmd.PersistentFlags().DurationVar(&statTick, "statictics-tick", tick, "duration for internal statistics agregation (Default: same as --tick)")

}

func StartServer(tsdb interfaces.TSDB, command *cobra.Command, args []string) (err error) {

	tStart := time.Now().UnixNano()
	pregenerated := tsdb.GenerateMetrics(maxMetrics)
	log.SLog.Infow("Metrics generated", "timer_ns", int((time.Now().UnixNano()-tStart)/1000000))

	log.SLog.Debug("Trying to generate senders")
	senders, err := generateSenders(tsdb, args)
	if err != nil {
		log.Log.Error(err.Error())
		return err
	}
	log.SLog.Infof("Created %d senders for %d hosts...", len(senders), len(args))
	server.StartServer(pregenerated, senders, startCount, tick, statTick, ListenURL, statDisable)
	return nil
}

func generateSenders(tsdb interfaces.TSDB, args []string) ([]interfaces.Sender, error) {
	array := []interfaces.Sender{}
	if len(args) == 0 {
		return array, fmt.Errorf("Please, add at least 1 destenation host")
	}
	for _, senderString := range args {
		uri, err := url.Parse(senderString)
		if err != nil {
			return array, err
		}
		for n := 0; n < parallel; n++ {
			sender, err := tsdb.NewSender(uri)
			if err != nil {
				return array, err
			}
			array = append(array, sender)
		}
	}
	return array, nil
}
