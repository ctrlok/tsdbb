package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/ctrlok/tsdbb/server"
	"github.com/spf13/cobra"
)

var startCount int
var parallel int
var tick time.Duration
var maxMetrics int
var statDisable bool
var listenURL string
var statTick time.Duration

// benchCmd represents the bench command
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	RootCmd.AddCommand(benchCmd)
	benchCmd.PersistentFlags().IntVarP(&startCount, "start-count", "c", 1000, "Start count of metrics which will be send")
	benchCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 6, "count of workers for each server in args")
	benchCmd.PersistentFlags().DurationVarP(&tick, "tick", "t", 1*time.Second, "retention period")
	benchCmd.PersistentFlags().IntVar(&maxMetrics, "maximum-metrics", 100000000, "maximum of metrics, which would be sended for one tick")
	benchCmd.PersistentFlags().BoolVar(&statDisable, "statistics-disable", false, "disable internal metrics")
	benchCmd.PersistentFlags().DurationVar(&statTick, "statictics-tick", tick, "duration for internal statistics agregation (Default: same as --tick)")
	benchCmd.PersistentFlags().StringVarP(&listenURL, "listen", "l", "127.0.0.1:8080", "set host:port for listening. Examples: 9090, :9090, 127.0.0.1:9090, 0.0.0.0:80")

}

func startServer(tsdb interfaces.TSDB, cmd *cobra.Command, args []string) (err error) {

	tStart := time.Now().UnixNano()
	pregenerated := tsdb.GenerateMetrics(maxMetrics)
	server.Logger.Info("Metrics generated", zap.Int("timer_ns", int((time.Now().UnixNano()-tStart)/1000000)))
	senders, err := generateSenders(tsdb, args)
	if err != nil {
		server.Logger.Error(err.Error())
		return err
	}
	server.StartServer(pregenerated, senders, startCount, tick, statTick, listenURL, statDisable)
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

func parseListen(s string) string {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s
	}
	return fmt.Sprint(":", s)
}
