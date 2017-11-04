package main

import (
	"github.com/jagandecapri/vision/orunada/utils"
	"github.com/jagandecapri/vision/orunada/preprocess"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sort"
	"time"
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/jagandecapri/vision/orunada/server"
	//"os"
)

func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort", "nbSrcs", "nbDsts", "perSyn", "perAck", "perRST", "perFIN", "perCWR", "perURG", "avgPktSize", "meanTTL")
	sort.Strings(sorter)
	return sorter
}

type DimMinMax struct{
	Min, Max float64
 	Range float64
}

func scale(elem float64, scale_factor float64) float64{
	return elem * scale_factor
}

func norm_mat(elem float64, col_min float64, col_max float64) float64{
	return (elem - col_min)/(col_max - col_min)
}

func normalize(mat []tree.Point, sorter []string) ([]tree.Point, map[string]DimMinMax){
	rows := len(mat)

	dim_min_max := map[string]DimMinMax{}
	for _,c := range sorter {
		min := mat[0].Vec_map[c]
		max := mat[0].Vec_map[c]
		for j := 0; j < rows; j++{
			val := mat[j].Vec_map[c]
			if val < min{
				min = val
			} else if  val > max{
				max = val
			}
		}
		range_ := max - min
		dim_min_max[c] = DimMinMax{min, max, range_}
	}

	for i := 0; i < rows; i++{
		for _, c := range sorter{
			col_min := dim_min_max[c].Min
			col_max := dim_min_max[c].Max
			elem := mat[i].Vec_map[c]
			if col_min == 0 && col_max == 0{
				//mat[i].Vec_map[c] = scale(elem, float64(scale_factor))
				mat[i].Vec_map[c] = elem
			} else {
				//mat[i].Vec_map[c] = scale(norm_mat(elem, col_min, col_max), float64(scale_factor)) //(elem - col_min)/(col_max - col_min)
				mat[i].Vec_map[c] = norm_mat(elem, col_min, col_max) //(elem - col_min)/(col_max - col_min)
			}
		}
	}

	//Assign normalized min-max
	for _,c := range sorter{
		norm_col_min := 0.0
		tmp := dim_min_max[c]
		tmp.Min = norm_col_min
		norm_col_max := 1.0
		tmp.Max = norm_col_max
		dim_min_max[c] = tmp
	}

	return mat, dim_min_max
}

//TODO: Refactor this as this is duplicate with window_time_slide.go
var delta_t time.Duration = 300 * time.Millisecond
var window time.Duration = 15 * time.Second
var window_arr_len = int(window.Seconds()/delta_t.Seconds())
var time_counter time.Time
var scale_factor = 5

type Config struct{
	min_dense_points int
	min_cluster_points int
}

func updateFS(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspaces map[[2]string]preprocess.Subspace, config Config){
	base_matrix := []tree.Point{}
	point_ctr := 0
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x := packet_acc.ExtractDeltaPacketFeature()
			point_ctr += 1
			p := tree.Point{Id: point_ctr, Vec_map: x}
			if len(base_matrix) < window_arr_len - 1{
				base_matrix = append(base_matrix, p)
			} else if len(base_matrix) == window_arr_len - 1{
				base_matrix = append(base_matrix, p)
				//TODO: normalization need to be done here, x_old and x_new will be what here?
			} else {
				fmt.Println("flow processing")
				base_matrix = append(base_matrix, p)
				norm_mat, _  := normalize(base_matrix, sorter)
				for _, subspace := range subspaces{
					subspace.ComputeSubspace(norm_mat)
					subspace.Cluster(config.min_dense_points, config.min_cluster_points)
					//os.Exit(2)
				}
				base_matrix = base_matrix[1:]
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
	subspaces := make(map[[2]string]preprocess.Subspace)

	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		Int_tree := tree.NewIntervalTree(uint64(dim))
		Units := tree.NewUnits()
		subspace := preprocess.Subspace{Interval_tree: &Int_tree, Units: &Units}
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
