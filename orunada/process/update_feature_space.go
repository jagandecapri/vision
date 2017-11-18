package process

import (
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/jagandecapri/vision/orunada/server"
	"github.com/jagandecapri/vision/orunada/tree"
	"sync"
	"fmt"
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
					done := make(chan struct{})
					c := ParallelClustering(done, subspaces, config, x_old, x_new_update)
					for r := range c{
						fmt.Printf("%v", r)
					}
					close(done)
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

func ParallelClustering(done chan struct{}, subspaces map[[2]string]Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point)(<- chan result){
	c := make(chan result)
	go func(){
		var wg sync.WaitGroup

		for _, subspace := range subspaces{

			wg.Add(1)
			go func(subspace Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point){
				subspace.ComputeSubspace(x_old, x_new_update)
				subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
				select {
				case c <- result{}:
				case <- done:
				}
				wg.Done()
			}(subspace, config, x_old, x_new_update)

			//select {
			//case <- done:
			//	return errors.New("Canceled")
			//default:
			//	return nil
			//}
		}

		go func(){
			wg.Wait()
			close(c)
		}()
	}()

	return c
}