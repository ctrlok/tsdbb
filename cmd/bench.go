// Copyright © 2017 Vsevolod Poliakov <ctrlok@gmail.com>
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

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ctrlok/tsdbb/log"
	"github.com/spf13/cobra"
)

// benchCmd represents the bench command
var benchCmd = &cobra.Command{
	Use:              "bench",
	Short:            "Start benchmark server",
	PersistentPreRun: preRunBench,
}

func init() {
	RootCmd.AddCommand(benchCmd)
	benchCmd.PersistentFlags().IntVarP(&Options.StartCount, "count", "c", 1000, "Start count of metrics which will be send")
	benchCmd.PersistentFlags().IntVar(&Options.StartStep, "step", 100, "default value for step")
	benchCmd.PersistentFlags().IntVarP(&Options.Parallel, "parallel", "p", 2, "Count of parallel workers for each server")
	benchCmd.PersistentFlags().DurationVarP(&Options.Tick, "tick", "t", time.Second, "retention period")
	benchCmd.PersistentFlags().StringP("listen", "l", ":8080", "adress for internal http server")
}

func preRunBench(cmd *cobra.Command, args []string) {
	rootPreRun(cmd, args)
	log.SLogger.Debug("Start preRun for bench")
	listenFlag := cmd.Flag("listen")
	Options.ListenURL = listenFlag.Value.String()
	if os.Getenv("LISTEN") != "" && !listenFlag.Changed {
		log.SLogger.Debug("Get listen paramether from os env")
		Options.ListenURL = os.Getenv("LISTEN")
	}
	Options.ListenURL = parseListen(Options.ListenURL)
	log.SLogger.Debugf("Setting listen url to %s", Options.ListenURL)
}

func parseListen(s string) string {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s
	}
	return fmt.Sprint(":", s)
}
