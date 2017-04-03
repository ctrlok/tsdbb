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
import "time"

var startCount int
var parallel int
var tick time.Duration
var maxMetrics int
var statDisable bool

// benchCmd represents the bench command
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	RootCmd.AddCommand(benchCmd)
	benchCmd.PersistentFlags().IntVarP(&startCount, "start-count", "c", 1000, "Start count of metrics which will be send")
	benchCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 6, "How much workers you should start to send metrics")
	benchCmd.PersistentFlags().DurationVarP(&tick, "tick", "t", 10*time.Second, "retention period")
	benchCmd.PersistentFlags().IntVar(&maxMetrics, "maximum-metrics", 10000000, "maximum of metrics, which would be sended for one tick")
	benchCmd.PersistentFlags().BoolVar(&statDisable, "statistics-disable", false, "disable internal metrics")

}
