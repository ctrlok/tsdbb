// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import "github.com/spf13/cobra"

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "A brief description of your command",
}

func init() {
	benchCmd.AddCommand(writeCmd)
	writeCmd.PersistentFlags().Bool("no-stats", false, "Disable internal statistics")
	writeCmd.PersistentFlags().String("statsd", "udp://localhost:8125", "Set statsd adress")
	writeCmd.PersistentFlags().String("graphite", "tcp://localhost:2003", "Set graphite adress")
}

func writePreRun(cmd *cobra.Command, args []string) {
	rootPreRun(cmd, args)

}
