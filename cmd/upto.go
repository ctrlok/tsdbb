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
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ctrlok/tsdbb/log"
	"github.com/ctrlok/tsdbb/server"
	"github.com/spf13/cobra"
)

var step int
var uris []string

// uptoCmd represents the upto command
var uptoCmd = &cobra.Command{
	Use:   "upto",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			log.SLogger.Fatal(err)
		}
		msg := server.ControlMessages{Count: count, Step: step}
		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			log.SLogger.Fatal(err)
		}
		for _, uri := range uris {
			log.SLogger.Debug("Send message to ", uri)
			log.SLogger.Debugf("Message is: %s", jsonMessage)
			i, err := http.Post(uri+"/upto", "application/json", bytes.NewBuffer(jsonMessage))
			if err != nil {
				log.SLogger.Fatal(err)
			}
			if i.StatusCode != http.StatusOK {
				log.SLogger.Fatalf("Server response with status %v", i.StatusCode)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(uptoCmd)

	uptoCmd.PersistentFlags().IntVarP(&step, "step", "s", 0, "Set new step")
	uptoCmd.PersistentFlags().StringArrayVarP(&uris, "servers", "l",
		[]string{"http://localhost:8080"}, "Set dst servers by comma")
}
