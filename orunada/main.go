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
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/jagandecapri/vision/orunada/tree"
	//"os"
	"github.com/jagandecapri/vision/orunada/server"
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

func getSubspace(subspace_key []string, mat []tree.Point, interval_tree augmentedtree.Tree) []tree.PointContainer {
	pnt_containers := []tree.PointContainer{}
	for _, p := range mat{
		tmp := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		tmp1 := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		int_container := tree.IntervalContainer{Id: 1, Range: tree.Range{Low: tmp, High: tmp1}, Scale_factor: scale_factor}
		interval := interval_tree.Query(int_container)
		if len(interval) > 0{
			Vec := []float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
			pnt_container := tree.PointContainer{
				Unit_id: int(interval[0].ID()),
				Vec: Vec,
				Point: p,
			}
			pnt_containers = append(pnt_containers, pnt_container)
			fmt.Println("Interval found", int_container, interval)
			//os.Exit(2)
		} else {
			fmt.Println("Empty interval:", int_container, interval)
		}
	}
	return pnt_containers
}

func updateFS(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspace_keys [][]string,
interval_trees map[[2]string]augmentedtree.Tree, units map[[2]string]tree.Units){
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
				for _, subspace_key := range subspace_keys{
					keys := [2]string{}
					copy(keys[:], subspace_key)
					getSubspace(subspace_key, norm_mat, interval_trees[keys])
					//subspace := getSubspace(subspace_key, norm_mat, interval_trees[keys])
					//tmp := kd_tree[keys]
					//for _, point := range subspace{
					//	tmp.AddToStore(point.Unit_id, point)
					//}
					//Clustering(tmp, 50, 1000)
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
	int_trees := make(map[[2]string]augmentedtree.Tree)
	units_arr := make(map[[2]string]tree.Units)
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
	intervals := tree.IntervalBuilder(ranges, scale_factor)
	units := tree.UnitsBuilder(ranges, dim)

	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		int_trees[tmp] = tree.NewIntervalTree(uint64(dim))
		units_arr[tmp] = tree.NewUnits()

		for _, interval := range intervals{
			int_trees[tmp].Add(interval)
		}

		for rg, unit := range units{
			units := units_arr[tmp]
			units.AddUnit(&unit, rg)
		}
	}

	//os.Exit(2)
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)
	go preprocess.WindowTimeSlide(ch, acc, quit)
	go updateFS(acc, data, sorter, subspace_keys, int_trees, units_arr)

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
