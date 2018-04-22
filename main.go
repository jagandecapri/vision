package main

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/process"
	"log"
	"sort"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/jagandecapri/vision/anomalies"
	"github.com/jagandecapri/vision/utils"
	"github.com/jagandecapri/vision/data"
	"github.com/jagandecapri/vision/cmd"
)


func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort",
		"nbSrcs", "nbDsts", "perSYN", "perACK", "perRST", "perFIN",
			"perCWR", "perURG", "perICMP", "avgPktSize", "meanTTL")
	sort.Strings(sorter)
	return sorter
}

func main(){
	cmd.Execute()

	if cmd.PrepareData{
		data.Run(cmd.PcapFilePath, cmd.DbNamePrepareData, cmd.DeltaTPrepareData)
	}

	if cmd.ClusterData{
		log.SetOutput(&lumberjack.Logger{
			Filename:   cmd.LogPath,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
			Compress:   true, // disabled by default
		})

		done := make(chan struct{})
		sorter:= getSorter()
		config := utils.Config{Min_dense_points: cmd.MinDensePoints, Min_cluster_points: cmd.MinClusterPoints, Num_cpu: cmd.NumCpu}
		subspace_channel_containers := anomalies.ClusteringBuilder(config, done)

		acc_c_receive := preprocess.AccumulatorChannels{
			AggSrc: make(preprocess.AccumulatorChannel),
			AggDst: make(preprocess.AccumulatorChannel),
			AggSrcDst: make(preprocess.AccumulatorChannel),
		}

		sql := data.NewSQLRead(cmd.DbNameClusterData, cmd.DeltaTClusterData)
		sql.ReadFromDb(acc_c_receive)

		acc_c_send := process.UpdateFeatureSpaceBuilder(subspace_channel_containers, sorter)
		preprocess.WindowTimeSlideSimulator(acc_c_receive, acc_c_send, cmd.DeltaTClusterData)
		<-done
		// ch := make(chan preprocess.PacketData)
		// BootServer(http_data)
	}



}
