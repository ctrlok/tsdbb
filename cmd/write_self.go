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
	"context"

	"github.com/ctrlok/tsdbb/interfaces/self"
	"github.com/ctrlok/tsdbb/log"
	"github.com/ctrlok/tsdbb/server"
	"github.com/spf13/cobra"
)

// graphiteCmd represents the graphite command
var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ctx1 := context.WithValue(context.Background(), log.KeyBenchType, "write")
		ctx2 := context.WithValue(ctx1, log.KeyTSDBType, "self")
		basic := &self.Basic{}
		Options.Servers = []string{""}
		server.StartServer(basic, Options, ctx2)
	},
}

func init() {
	writeCmd.AddCommand(selfCmd)
}
