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

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ctrlok/tsdbb/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ListenURL string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "tsdb-bench",
	Short:            "A brief description of your application a",
	Long:             ``,
	PersistentPreRun: PreRun,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&ListenURL, "listen-url", "l", "127.0.0.1:8080", "set host:port for listening. Examples: 9090, :9090, 127.0.0.1:9090, 0.0.0.0:80")
	RootCmd.PersistentFlags().BoolP("debug", "D", false, "set log level to debug")
	RootCmd.PersistentFlags().Bool("json", false, "Save logs in json format")
	RootCmd.AddCommand(BenchCmd)
}

func PreRun(cmd *cobra.Command, args []string) {
	listenFlag := cmd.Flag("listen-url")
	if os.Getenv("LISTEN") != "" && !listenFlag.Changed {
		ListenURL = os.Getenv("LISTEN")
	}
	ListenURL = parseListen(ListenURL)

	var config zap.Config

	if debugFlag, _ := cmd.Flags().GetBool("debug"); debugFlag {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		config.Encoding = "json"
		config.EncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	log.Log, _ = config.Build()
	log.SLog = log.Log.Sugar()

	log.Log.Debug("Log was seted")

}

func parseListen(s string) string {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s
	}
	return fmt.Sprint(":", s)
}
