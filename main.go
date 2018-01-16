package main

import (
	"github.com/jagandecapri/vision/utils"
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/process"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sort"
	"github.com/jagandecapri/vision/tree"
	"github.com/jagandecapri/vision/server"
	"flag"
)

var scale_factor = 5

func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort", "nbSrcs", "nbDsts", "perSyn", "perAck", "perRST", "perFIN", "perCWR", "perURG", "avgPktSize", "meanTTL")
	sort.Strings(sorter)
	return sorter
}

func main(){
	num_cpu := flag.Int("num-cpu", 0, "Number of CPUs to use")
	flag.Parse()

	data := make(chan server.HttpData)
	go BootServer(data)

	sorter:= getSorter()
	subspace_keys := utils.GetKeyComb(sorter, 2)
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	subspaces := make(map[[2]string]tree.Subspace)

	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		Int_tree := tree.NewIntervalTree(uint64(dim))
		grid := tree.NewGrid()
		ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
		intervals := tree.IntervalBuilder(ranges, scale_factor)
		units := tree.UnitsBuilder(ranges, dim)

		subspace := tree.Subspace{Grid: &grid, Subspace_key: tmp, Scale_factor: scale_factor}
		subspace.SetIntervalTree(&Int_tree)
		for _, interval := range intervals{
			Int_tree.Add(interval)
		}

		for _, unit := range units{
			grid.AddUnit(&unit)
		}
		grid.SetupGrid(interval_length)
		subspaces[tmp] = subspace
	}

	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)

	config := process.Config{Min_dense_points: 10, Min_cluster_points: 15, Execution_type: process.PARALLEL, Num_cpu: *num_cpu}
	go preprocess.WindowTimeSlide(ch, acc, quit)
	go process.UpdateFeatureSpace(acc, data, sorter, subspaces, config)

	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			quit <- 0
			log.Fatal(err)
		} else if err == io.EOF {
			quit <- 0
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- preprocess.PacketData{packet, ci}
		}
	}

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
