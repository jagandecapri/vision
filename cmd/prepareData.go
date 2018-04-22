// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"time"
)

// prepareDataCmd represents the prepareData command
var prepareDataCmd = &cobra.Command{
	Use:   "prepareData",
	Short: "Prepare network flow data",
	Long: `Prepare network flow data`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prepareData called")
		PrepareData = true
	},
}

var PrepareData bool
var PcapFilePath string
var DbNamePrepareData string
var DeltaTPrepareData time.Duration

func init() {
	rootCmd.AddCommand(prepareDataCmd)

	prepareDataCmd.Flags().StringVarP(&PcapFilePath, "pcap-file-path", "", "", "pcap-file-path")
	prepareDataCmd.MarkFlagRequired("pcap-file-path")
	prepareDataCmd.Flags().StringVarP(&DbNamePrepareData, "db-name", "", "", "Name of SQLite database used to store network flow data")
	prepareDataCmd.MarkFlagRequired("db-name")
	prepareDataCmd.Flags().DurationVarP(&DeltaTPrepareData, "delta-t", "", 300 * time.Millisecond, "Delta time")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prepareDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prepareDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
