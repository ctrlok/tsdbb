package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces/graphite"
	"github.com/spf13/cobra"
)

// graphiteCmd represents the graphite command
var graphiteCmd = &cobra.Command{
	Use:   "graphite [servers] [flags]",
	Short: "bench graphite servers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tsdb := &graphite.TSDB{}
		prefix, _ := cmd.Flags().GetString("prefix")
		tsdb.GeneratorPrefix = []byte(prefix)
		startServer(tsdb, cmd, args)
	},
}

func init() {
	benchCmd.AddCommand(graphiteCmd)
	graphiteCmd.PersistentFlags().String("prefix", "metric", "it is prefix for graphite metrics")
}
