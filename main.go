package main

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/process"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sort"
	"github.com/jagandecapri/vision/server"
	"flag"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/jagandecapri/vision/anomalies"
	"github.com/jagandecapri/vision/utils"
	"os"
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
	log.SetOutput(&lumberjack.Logger{
		Filename:   "C:\\Users\\Jack\\go\\src\\github.com\\jagandecapri\\vision\\logs\\lumber_log.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true, // disabled by default
	})

	num_cpu := flag.Int("num-cpu", 0, "Number of CPUs to use")
	flag.Parse()

	data := make(chan server.HttpData)
	ch := make(chan preprocess.PacketData)
	done := make(chan struct{})

	sorter:= getSorter()
	config := utils.Config{Min_dense_points: 10, Min_cluster_points: 15, Execution_type: utils.PARALLEL, Num_cpu: *num_cpu}

	BootServer(data)
	subspace_channel_containers := anomalies.ClusteringBuilder(config, done)
	accumulator_channels := process.UpdateFeatureSpaceBuilder(subspace_channel_containers, sorter, done)
	preprocess.WindowTimeSlide(ch, accumulator_channels, done)

	pcap_file_path := os.Getenv("PCAP_FILE")
	if pcap_file_path == ""{
		pcap_file_path = "C:\\Users\\Jack\\Downloads\\201705021400.pcap"
	}
	handleRead, err := pcap.OpenOffline(pcap_file_path)

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			close(done)
			log.Fatal(err)
		} else if err == io.EOF {
			close(done)
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- preprocess.PacketData{Data: packet, Metadata: ci}
		}
	}

	////HACK
	//subspace_keys := [][]string{{"nbSrcs", "avgPktSize"}, {"perICMP", "perSYN"}, {"nbSrcPort", "perICMP"}}
	//
	//min_interval := 0.0
	//max_interval := 1.0
	//interval_length := 0.1
	//dim := 2
	//subspaces := make(map[[2]string]tree.Subspace)
	//
	//for _, subspace_key := range subspace_keys{
	//	tmp := [2]string{}
	//	copy(tmp[:], subspace_key)
	//	Int_tree := tree.NewIntervalTree(uint64(dim))
	//	grid := tree.NewGrid()
	//	ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
	//	intervals := tree.IntervalBuilder(ranges, scale_factor)
	//	units := tree.UnitsBuilder(ranges, dim)
	//
	//	subspace := tree.Subspace{Grid: &grid, Subspace_key: tmp, Scale_factor: scale_factor}
	//	subspace.SetIntervalTree(&Int_tree)
	//	for _, interval := range intervals{
	//		Int_tree.Add(interval)
	//	}
	//
	//	for _, unit := range units{
	//		grid.AddUnit(&unit)
	//	}
	//	grid.SetupGrid(interval_length)
	//	subspaces[tmp] = subspace
	//}
	//
	//ch := make(chan preprocess.PacketData)
	//acc_c := make(chan preprocess.X_micro_slot)
	////quit := make(chan int)
	//
	//config := process.Config{Min_dense_points: 10, Min_cluster_points: 15, Execution_type: process.PARALLEL, Num_cpu: *num_cpu}
	//go preprocess.WindowTimeSlide(ch, acc_c)
	//go process.UpdateFeatureSpace(acc_c, data, sorter, subspaces, config)
	//
	//handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	//
	//if(err != nil){
	//	log.Fatal(err)
	//}
	//
	//for {
	//	data, ci, err := handleRead.ReadPacketData()
	//	if err != nil && err != io.EOF {
	//		close(ch)
	//		//quit <- 0
	//		log.Fatal(err)
	//	} else if err == io.EOF {
	//		close(ch)
	//		//quit <- 0
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
