package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces/graphite"
	"github.com/spf13/cobra"
)

// SelfCmd represents the self command
var selfGraphiteCmd = &cobra.Command{
	Use:   "self",
	Short: "testing internal speed",
	Long: `That method is primary for uderstanding how mush metrics you can send
without locking on syscals and actual sending date.`,
	Run: func(cmd *cobra.Command, args []string) {
		tsdb := &graphite.TSDB{DevNull: true}
		StartServer(tsdb, cmd, []string{""})
	},
}

func init() {
	GraphiteCmd.AddCommand(selfGraphiteCmd)
}
