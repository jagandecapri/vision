package process

import (
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/jagandecapri/vision/orunada/server"
	"github.com/jagandecapri/vision/orunada/tree"
	"sync"
	"fmt"
	"runtime"
	"time"
	"github.com/jagandecapri/vision/orunada/utils"
)

func UpdateFeatureSpace(acc chan preprocess.PacketAcc, data chan server.HttpData, sorter []string, subspaces map[[2]string]Subspace, config Config){
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
				norm_mat := preprocess.Normalize(base_matrix, sorter)
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

				if (config.Execution_type == SEQUENTIAL){
					c := SequentialClustering(subspaces, config, x_old, x_new_update)
					for r := range c{
						fmt.Println(r)
					}
				} else if (config.Execution_type == PARALLEL){
					//cur_CPU := runtime.GOMAXPROCS(1)
					cur_CPU := runtime.GOMAXPROCS(0) //gets the current number of cores
					num_clusterer := cur_CPU
					fmt.Println("cur num CPU", cur_CPU)
					func (){
						defer utils.TimeTrack(time.Now(), "Clustering")
						m := ParallelClustering(num_clusterer, subspaces, config, x_old, x_new_update)
						for _, r := range m{
							fmt.Printf("%v", r)
						}
						return
					}()
				}

				//os.Exit(2)
			}
		}
	}
}

type result struct{
	data int
}

func SequentialClustering(subspaces map[[2]string]Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point) []result{
	for _, subspace := range subspaces{
		subspace.ComputeSubspace(x_old, x_new_update)
		subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
	}
	return []result{}
}

type processPackage struct{
	subspace Subspace
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
		case c <- result{}:
		case <- done:
			return
		}
	}
}

func SubspaceIterator(done <- chan struct{}, subspaces map[[2]string]Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point) (<-chan processPackage){
	processPackages := make(chan processPackage)
	go func(){
		fmt.Println("Subspace len:", len(subspaces))
		defer close(processPackages)

		for _, subspace := range subspaces {
			processPackage := processPackage{
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

func ParallelClustering(num_clusterers int, subspaces map[[2]string]Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point) []result{
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

	m := []result{}
	for r := range c{
		//to extract needed data for server visualization
		m = append(m, r)
	}

	return m
}