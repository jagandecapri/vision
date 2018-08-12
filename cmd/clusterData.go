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
	"github.com/pkg/errors"
	"path/filepath"
	"os"
)

// clusterDataCmd represents the clusterData command
var clusterDataCmd = &cobra.Command{
	Use:   "clusterData",
	Short: "Cluster network flow data to detect anomalies",
	Long: `Cluster network flow data to detect anomalies.
Anomalies are found in log file.`,
	Args: func(cmd *cobra.Command, args []string) error {
		log_path, err := cmd.Flags().GetString("log-path")
		if err != nil{
			dir, _ := filepath.Split(log_path)

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return errors.New("Log path does not exist")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clusterData called")
		ClusterData = true
	},
}

var ClusterData bool
var DbNameClusterData string
var LogPath string
var NumCpu int
var MinDensePoints int
var MinClusterPoints int
var PointsMode string
var DeltaTClusterData time.Duration
var WindowArrayLen int
var NumKneeFlatPoints int
var KneeSmoothingWindow int
var KneeFindElbow bool
var CpuProfile bool
var MemProfile bool

func init() {
	rootCmd.AddCommand(clusterDataCmd)

	clusterDataCmd.Flags().StringVarP(&DbNameClusterData, "db-name", "","", "Name of SQLite database used for clustering")
	clusterDataCmd.MarkFlagRequired("db-name")
	clusterDataCmd.Flags().StringVarP(&LogPath, "log-path", "","", "Path to write log. Make sure path is valid.")
	clusterDataCmd.MarkFlagRequired("log-path")
	clusterDataCmd.Flags().IntVarP(&NumCpu, "num-cpu", "",0, "Number of CPUs to use")
	clusterDataCmd.Flags().IntVarP(&MinDensePoints, "min-dense-points", "", 10, "Minimum number of points to consider a unit as dense")
	clusterDataCmd.Flags().IntVarP(&MinClusterPoints, "min-cluster-points", "",10, "Minimum number of points to consider a group of units as cluster")
	clusterDataCmd.Flags().StringVarP(&PointsMode, "points-mode", "", "default", "Mode for point calculation in minimum number of points of a dense unit and cluster. Default mode is to take the number of points. Use 'percentage' to take the percentage of points.")
	clusterDataCmd.Flags().DurationVarP(&DeltaTClusterData, "delta-t", "", 300 * time.Millisecond, "Delta time")
	clusterDataCmd.Flags().IntVarP(&WindowArrayLen, "window-array-len", "", 2, "Number of micro slots for window time sliding")
	clusterDataCmd.Flags().IntVarP(&NumKneeFlatPoints, "num-knee-flat-points", "", 1, "number of 'flat' points that is required before considering a point as knee or elbow.")
	clusterDataCmd.Flags().IntVarP(&KneeSmoothingWindow, "knee-smoothing-window", "", 3, "The smootingWindow parameter is used to indicate the average used for the Gaussian kernel average smoother")
	clusterDataCmd.Flags().BoolVarP(&KneeFindElbow, "knee-find-elbow", "", false, "Indicates whether to find an elbow or a knee when the value of parameter is true or false respectively")
	clusterDataCmd.Flags().BoolVarP(&CpuProfile, "cpu-profile", "", false, "write cpu profile to this file")
	clusterDataCmd.Flags().BoolVarP(&MemProfile, "mem-profile", "", false, "write memory profile to this file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
