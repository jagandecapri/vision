package process

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/server"
	"github.com/jagandecapri/vision/tree"
	"sync"
	"fmt"
	"runtime"
	"time"
	"github.com/jagandecapri/vision/utils"
)


func UpdateFeatureSpace(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspaces map[[2]string]tree.Subspace, config Config){
	base_matrix := []tree.Point{}
	l := utils.Logger{}
	logger := l.New()
	point_ctr := 0
	count := 0
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
				norm_mat := preprocess.Normalize(base_matrix, sorter)
				var x_old, x_new_update []tree.Point

				if len(base_matrix) == preprocess.WINDOW_ARR_LEN{
					//fmt.Println("before flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					//fmt.Println("before flow processing")
					x_old, x_new_update = []tree.Point{}, norm_mat
				} else if len(base_matrix) > preprocess.WINDOW_ARR_LEN{
					//fmt.Println("flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					//fmt.Println("flow processing")
					x_old, x_new_update = []tree.Point{norm_mat[0]}, norm_mat[1:]
					base_matrix = base_matrix[1:]
				}

				m := map[[2]string]tree.Subspace{}
				if (config.Execution_type == SEQUENTIAL){
					m = SequentialClustering(subspaces, config, x_old, x_new_update)

				} else if (config.Execution_type == PARALLEL){
					num_clusterer := runtime.GOMAXPROCS(config.Num_cpu) //gets the current number of cores
					func (){
						defer utils.TimeTrack(time.Now(),  "Clustering", num_clusterer, logger)
						//ParallelClustering(num_clusterer, subspaces, config, x_old, x_new_update)
						m = ParallelClustering(num_clusterer, subspaces, config, x_old, x_new_update)
						count++
						return
					}()
				}

				//http_data := server.HttpData{}
				//for subspace_key, subspace := range m{
				//	subspace_key_join := strings.Join(subspace_key[:], "-")
				//	point_cluster := GetVisualizationData(subspace)
				//	http_data.Point_cluster[subspace_key_join] = point_cluster
				//}
				//http_data.Points = m[0].Grid.
				//fmt.Println("Send HTTP data: ", http_data)
				//data <- http_data
			}
		}
	}
}

func SequentialClustering(subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point)  map[[2]string]tree.Subspace{
	for _, subspace := range subspaces{
		subspace.ComputeSubspace(x_old, x_new_update)
		subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
	}
	return subspaces
}

type result struct{
	subspace_key [2]string
	subspace tree.Subspace
}

type processPackage struct{
	subspace_key [2]string
	subspace tree.Subspace
	config Config
	x_old []tree.Point
	x_new_update []tree.Point
}

func Clusterer(done <-chan struct{}, processPackages <-chan processPackage , c chan<- result){
	for processPackage := range processPackages{
		subspace := processPackage.subspace
		config := processPackage.config
		x_old := processPackage.x_old
		x_new_update := processPackage.x_new_update
		subspace.ComputeSubspace(x_old, x_new_update)
		subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
		select {
		case c <- result{processPackage.subspace_key, subspace}:
		case <- done:
			return
		}
	}
}

func SubspaceIterator(done <- chan struct{}, subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point) (<-chan processPackage){
	processPackages := make(chan processPackage)
	go func(){
		//fmt.Println("Subspace len:", len(subspaces))
		defer close(processPackages)

		for _, subspace := range subspaces {
			processPackage := processPackage{
				subspace_key: subspace.Subspace_key,
				subspace: subspace,
				config: config,
				x_old: x_old,
				x_new_update: x_new_update,
			}
			select {
			case processPackages <- processPackage:
			case <-done:
				break
			}
		}
		return
	}()
	return processPackages
}

func ParallelClustering(num_clusterers int, subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point)  map[[2]string]tree.Subspace{
	done := make(chan struct{})
	defer close(done)

	processPackages := SubspaceIterator(done, subspaces, config ,x_old, x_new_update)

	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(num_clusterers)
	for i := 0; i < num_clusterers; i++{
		go func(){
			Clusterer(done, processPackages, c)
			wg.Done()
		}()
	}
	go func(){
		wg.Wait()
		close(c)
	}()

	m := make(map[[2]string]tree.Subspace)
	for r := range c{
		//to extract needed data for server visualization
		m[r.subspace_key] = r.subspace
	}

	return m
}