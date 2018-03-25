package main

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/process"
	"log"
	"sort"
	"github.com/jagandecapri/vision/server"
	"flag"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/jagandecapri/vision/anomalies"
	"github.com/jagandecapri/vision/utils"
	"os"
	"github.com/jagandecapri/vision/data"
	"time"
)

var scale_factor = 5

func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort",
		"nbSrcs", "nbDsts", "perSYN", "perACK", "perRST", "perFIN",
			"perCWR", "perURG", "perICMP", "avgPktSize", "meanTTL")
	sort.Strings(sorter)
	return sorter
}

func main(){
	log_path := os.Getenv("LOG_FILE")
	if log_path == ""{
		log_path = "C:\\Users\\Jack\\go\\src\\github.com\\jagandecapri\\vision\\logs\\lumber_log.log"
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   log_path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true, // disabled by default
	})

	num_cpu := flag.Int("num-cpu", 0, "Number of CPUs to use")
	flag.Parse()

	http_data := make(chan server.HttpData)
	done := make(chan struct{})
	sorter:= getSorter()
	config := utils.Config{Min_dense_points: 10, Min_cluster_points: 15, Execution_type: utils.PARALLEL, Num_cpu: *num_cpu}

	delta_t := 300 * time.Millisecond

	BootServer(http_data)
	subspace_channel_containers := anomalies.ClusteringBuilder(config, done)

	acc_c_receive := preprocess.AccumulatorChannels{
		AggSrc: make(preprocess.AccumulatorChannel),
		AggDst: make(preprocess.AccumulatorChannel),
		AggSrcDst: make(preprocess.AccumulatorChannel),
	}

	sql := data.NewSQLRead("./201705021400.db", delta_t)
	sql.ReadFromDb(acc_c_receive)

	acc_c_send := process.UpdateFeatureSpaceBuilder(subspace_channel_containers, sorter)
	preprocess.WindowTimeSlideSimulator(acc_c_receive, acc_c_send, delta_t)
	<-done
	// ch := make(chan preprocess.PacketData)
	// BootServer(http_data)
	//subspace_channel_containers := anomalies.ClusteringBuilder(config, done)
	//accumulator_channels := process.UpdateFeatureSpaceBuilder(subspace_channel_containers, sorter, done)
	//preprocess.WindowTimeSlide(ch, accumulator_channels, done)
	//
	//pcap_file_path := os.Getenv("PCAP_FILE")
	//if pcap_file_path == ""{
	//	pcap_file_path = "C:\\Users\\Jack\\Downloads\\201705021400.pcap"
	//}
	//
	//handleRead, err := pcap.OpenOffline(pcap_file_path)
	//
	//if(err != nil){
	//	log.Fatal(err)
	//}
	//
	//for {
	//	data, ci, err := handleRead.ReadPacketData()
	//	if err != nil && err != io.EOF {
	//		close(done)
	//		log.Fatal(err)
	//	} else if err == io.EOF {
	//		close(done)
	//		break
	//	} else {
	//		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	//		ch <- preprocess.PacketData{Data: packet, Metadata: ci}
	//	}
	//}

	//if handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever); err != nil {
	//	panic(err)
	//} else {
	//	for {
	//		data, ci, err := handle.ReadPacketData()
	//		if err != nil && err != io.EOF {
	//			quit <- 0
	//			log.Fatal(err)
	//		} else if err == io.EOF {
	//			quit <- 0
	//			break
	//		} else {
	//			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	//			ch <- preprocess.PacketData{packet, ci}
	//		}
	//	}
	//}
}
