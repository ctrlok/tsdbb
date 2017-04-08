package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces/graphite"
	"github.com/spf13/cobra"
)

// selfCmd represents the self command
var selfGraphiteCmd = &cobra.Command{
	Use:   "self",
	Short: "testing internal speed",
	Long: `That method is primary for uderstanding how mush metrics you can send
without locking on syscals and actual sending date.`,
	Run: func(cmd *cobra.Command, args []string) {
		tsdb := &graphite.TSDB{DevNull: true}
		startServer(tsdb, cmd, []string{""})
	},
}

func init() {
	graphiteCmd.AddCommand(selfGraphiteCmd)
}
