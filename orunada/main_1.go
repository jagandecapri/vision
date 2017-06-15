package main

import (
	"github.com/jagandecapri/vision/orunada/utils"
	"github.com/jagandecapri/vision/orunada/grid"
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/cockroachdb/apd"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sort"
	//"os"
	"time"
	"os"
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

func normalize(mat []grid.Point, sorter []string) ([]grid.Point, map[string]DimMinMax){
	rows := len(mat)
	//cols := len(mat[0])

	dim_min_max := map[string]DimMinMax{}
	for _,c := range sorter {
		min := mat[0].Vec[c]
		max := mat[0].Vec[c]
		for j := 0; j < rows; j++{
			val := mat[j].Vec[c]
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
			elem := mat[i].Vec[c]
			if col_min == 0 && col_max == 0{
				mat[i].Norm_vec[c] = elem
			} else {
				mat[i].Norm_vec[c] = (elem - col_min)/(col_max - col_min)
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

func getSubspace(subspace_keys [][]string, mat []grid.Point) map[[2]string][]grid.Point{
	subspaces := map[[2]string][]grid.Point{}
	for _, subspace_k := range subspace_keys{
		key := [2]string{}
		copy(key[:], subspace_k)
		subspace := []grid.Point{}
		for _, p:= range mat{
			sub_point := grid.Point{Id: p.Id}
			tmp := make(map[string]float64)
			for i := 0; i < len(subspace_k); i++{
				key := subspace_k[i]
				tmp[key] = p.Norm_vec[key]
			}
			sub_point.Norm_vec = tmp
			sub_point.Sorter = key
			subspace = append(subspace, sub_point)
		}
		subspaces[key] = subspace
	}
	return subspaces
}

func updateFS(acc chan preprocess.PacketAcc, sorter []string, subspace_keys [][]string, grids map[[2]string]grid.Grid){
	base_matrix := []grid.Point{}
	point_ctr := 0
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x := packet_acc.ExtractDeltaPacketFeature()
			point_ctr += 1
			p := grid.Point{Id: point_ctr, Vec: x, Norm_vec: make(map[string]float64)}
			if len(base_matrix) < window_arr_len - 1{
				base_matrix = append(base_matrix, p)
			} else if len(base_matrix) == window_arr_len - 1{
				base_matrix = append(base_matrix, p)
				//TODO: normalization need to be done here, x_old and x_new will be what here?
			} else {
				base_matrix = append(base_matrix, p)
				norm_mat, _ := normalize(base_matrix, sorter)
				subspaces := getSubspace(subspace_keys, norm_mat)
				for key, subspace := range subspaces{
					g := grids[key]
					x_old, x_update, x_new := []grid.Point{subspace[0]}, subspace[1:len(subspace)-2], []grid.Point{subspace[len(subspace)-1]}
					g.Assign(x_old)
					g.Assign(x_update)
					g.Assign(x_new)
					os.Exit(2)
				}
				base_matrix = base_matrix[1:]
			}

		}
	}
}

func main(){
	sorter:= getSorter()
	grids := map[[2]string]grid.Grid{}
	subspace_keys := utils.GetKeyComb(sorter, 2)
	ctx := apd.BaseContext.WithPrecision(6)
	for _, subspace_key := range subspace_keys{
		tmp := [2]string{}
		copy(tmp[:], subspace_key)
		grid := grid.Grid{}
		grid.Build2DGrid(subspace_key, ctx)
		grids[tmp] = grid
	}

	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)
	go preprocess.WindowTimeSlide(ch, acc, quit)
	go updateFS(acc, sorter, subspace_keys, grids)

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
