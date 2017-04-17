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

	"go.uber.org/zap"

	"github.com/ctrlok/tsdbb/log"
	"github.com/ctrlok/tsdbb/server"
	"github.com/spf13/cobra"
)

var Options server.Options

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "tsdbb",
	Short:            "A brief description of your application",
	PersistentPreRun: rootPreRun,
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
	RootCmd.PersistentFlags().BoolVarP(&log.DebugLevel, "debug", "D", false, "Set log level to debug")
	RootCmd.PersistentFlags().Bool("json", false, "format log to json")
}

func rootPreRun(cmd *cobra.Command, args []string) {
	fmt.Println("LOG START!!!!")
	var config zap.Config

	if log.DebugLevel {
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
	log.Logger, _ = config.Build()
	log.SLogger = log.Logger.Sugar()
	log.SLogger.Debugw("Log initialized as", "type", config.Encoding, "debug", log.DebugLevel)
}
