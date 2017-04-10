package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces/graphite"
	"github.com/spf13/cobra"
)

// GraphiteCmd represents the graphite command
var GraphiteCmd = &cobra.Command{
	Use:   "graphite [servers] [flags]",
	Short: "bench graphite servers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tsdb := &graphite.TSDB{}
		prefix, _ := cmd.Flags().GetString("prefix")
		tsdb.GeneratorPrefix = []byte(prefix)
		StartServer(tsdb, cmd, args)
	},
}

func init() {
	GraphiteCmd.PersistentFlags().String("prefix", "metric", "it is prefix for graphite metrics")
	BenchCmd.AddCommand(GraphiteCmd)
}
