// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"os"
	"os/signal"
	"strconv"

	"github.com/gotoolkit/miner/container"
	"github.com/gotoolkit/miner/db"
	"github.com/spf13/cobra"
)

const DockerAPIMinVersion string = "1.26"

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if botToken == "" {
			cmd.Usage()
			os.Exit(1)
		}

		db.Init()
		defer db.Close()

		os.Setenv("DOCKER_HOST", dockerHost)
		os.Setenv("DOCKER_TLS_VERIFY", strconv.FormatBool(tlsverify))
		os.Setenv("DOCKER_API_VERSION", DockerAPIMinVersion)

		client := container.NewClient(botToken)
		go client.StartWebHook()
		go client.StartBot(userID)
		// go client.StartInlineQuery()
		handleSignals()
	},
}

var (
	dockerHost string
	tlsverify  bool
	botToken   string
	userID     int
)

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serverCmd.Flags().IntVarP(&userID, "id", "I", 1, "Telegram author user")
	serverCmd.Flags().StringVarP(&botToken, "token", "T", "", "Telegram Bot Token")
	serverCmd.Flags().StringVarP(&dockerHost, "host", "H", "unix:///var/run/docker.sock", "Docker host to connect to")

}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
