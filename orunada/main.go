package main

import (
	"github.com/jagandecapri/vision/orunada/utils"
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/jagandecapri/vision/orunada/process"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sort"
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/jagandecapri/vision/orunada/server"
)

var scale_factor = 5

type Config struct{
	min_dense_points int
	min_cluster_points int
}

func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort", "nbSrcs", "nbDsts", "perSyn", "perAck", "perRST", "perFIN", "perCWR", "perURG", "avgPktSize", "meanTTL")
	sort.Strings(sorter)
	return sorter
}

func updateFS(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspaces map[[2]string]process.Subspace, config Config){
	base_matrix := []tree.Point{}
	point_ctr := 0
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x := packet_acc.ExtractDeltaPacketFeature()
			point_ctr += 1
			p := tree.Point{Id: point_ctr, Vec_map: x}
			if len(base_matrix) < preprocess.WINDOW_ARR_LEN - 1{
				base_matrix = append(base_matrix, p)
			} else {
				base_matrix = append(base_matrix, p)
				norm_mat, _  := preprocess.Normalize(base_matrix, sorter)
				var x_old, x_new_update []tree.Point

				if len(base_matrix) == preprocess.WINDOW_ARR_LEN{
					fmt.Println("before flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					fmt.Println("before flow processing")
					x_old, x_new_update = []tree.Point{}, norm_mat
				} else if len(base_matrix) > preprocess.WINDOW_ARR_LEN{
					fmt.Println("flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					fmt.Println("flow processing")
					x_old, x_new_update = []tree.Point{norm_mat[0]}, norm_mat[1:]
					base_matrix = base_matrix[1:]
				}
				for _, subspace := range subspaces{
					subspace.ComputeSubspace(x_old, x_new_update)
					subspace.Cluster(config.min_dense_points, config.min_cluster_points)
				}
				//os.Exit(2)
			}
		}
	}
}

func main(){
	data := make(chan server.HttpData)
	//go BootServer(data)

	sorter:= getSorter()
	subspace_keys := utils.GetKeyComb(sorter, 2)
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
	intervals := tree.IntervalBuilder(ranges, scale_factor)
	units := tree.UnitsBuilder(ranges, dim)
	subspaces := make(map[[2]string]process.Subspace)

	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		Int_tree := tree.NewIntervalTree(uint64(dim))
		Units := tree.NewUnits()
		subspace := process.Subspace{Interval_tree: &Int_tree, Units: &Units, Subspace_key: tmp, Scale_factor: scale_factor}
		for _, interval := range intervals{
			Int_tree.Add(interval)
		}

		for rg, unit := range units{
			Units.AddUnit(&unit, rg)
		}
		Units.SetupGrid(interval_length)
		subspaces[tmp] = subspace
	}

	//os.Exit(2)
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)

	config := Config{min_dense_points: 10, min_cluster_points: 15}
	go preprocess.WindowTimeSlide(ch, acc, quit)
	go updateFS(acc, data, sorter, subspaces, config)

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
}
