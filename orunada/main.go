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
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"github.com/jagandecapri/vision/orunada/tree"
	"os"
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
		tmp := []float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[0]]}
		tmp1 := []float64{p.Vec_map[subspace_key[1]], p.Vec_map[subspace_key[1]]}
		int_container := tree.IntervalContainer{Id: 1, Low: tmp, High: tmp1, Scale_factor: scale_factor}
		interval := interval_tree.Query(int_container)
		if len(interval) > 0{
			Vec := []float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
			pnt_container := tree.PointContainer{
				Unit_id: int(interval[0].ID()),
				Vec: Vec,
				Point: p,
			}
			pnt_containers = append(pnt_containers, pnt_container)
		} else {
			fmt.Println("Empty interval:", int_container, interval)
		}
	}
	return pnt_containers
}

func Clustering(kd_ext tree.KDTree_Extend, density int, distance int){
	out := make(chan tree.PointInterface)
	kd_ext.KDTree.BFSTraverseChan(out)
	for p := range out{
		fmt.Println(p)
	}
}

func updateFS(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspace_keys [][]string,
interval_trees map[[2]string]augmentedtree.Tree, kd_tree map[[2]string]tree.KDTree_Extend){
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
					os.Exit(2)
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
	int_trees := map[[2]string]augmentedtree.Tree{}
	kd_trees_ext := map[[2]string]tree.KDTree_Extend{}
	intervals := tree.IntervalBuilder(0.0, 1.0, 0.1, scale_factor)
	fmt.Println(intervals)

	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		int_tree := tree.NewIntervalTree(2)
		kd_tree_ext := tree.KDTree_Extend{&tree.KDTree{}, make(map[int]*tree.Unit)}
		int_trees[tmp] = int_tree
		kd_trees_ext[tmp] = kd_tree_ext

		for _, v := range intervals{
			unit_center := tree.PointContainer{Unit_id: v.Id, Vec: v.GetCenter()}
			unit := tree.Unit{Id: v.Id, IntervalContainer: &v, Dimension: 2, Center: unit_center}
			int_tree.Add(unit)
			kd_tree_ext.AddUnit(&unit)
		}
		break
	}

	//os.Exit(2)
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)
	go preprocess.WindowTimeSlide(ch, acc, quit)
	go updateFS(acc, data, sorter, subspace_keys, int_trees, kd_trees_ext)

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
