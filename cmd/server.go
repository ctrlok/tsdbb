package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/spf13/cobra"
)

func startServer(tsdb interfaces.TSDB, cmd *cobra.Command, args []string) (err error) {
	tsdb.GenerateMetrics(maxMetrics)
	// senders := generateSenders(tsdb, args)
	return nil
}

func generateSenders(tsdb interfaces.TSDB, args []string) []interfaces.Sender {
	s := tsdb.NewSender()
	return []interfaces.Sender{s}
}
