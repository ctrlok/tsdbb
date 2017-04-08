package cmd

import (
	"github.com/ctrlok/tsdbb/interfaces/self"
	"github.com/spf13/cobra"
)

// selfCmd represents the self command
var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "testing internal speed",
	Long: `That method is primary for uderstanding how mush metrics you can send without locking
on syscals and actual sending date. For understanding how much you can send metric to
provider, please use:

tsdbb bench graphite self [flags]`,
	Run: func(cmd *cobra.Command, args []string) {
		tsdb := &self.TSDB{}
		startServer(tsdb, cmd, []string{""})
	},
}

func init() {
	benchCmd.AddCommand(selfCmd)
}
