package process

import (
	"fmt"
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/jagandecapri/vision/orunada/server"
	"github.com/jagandecapri/vision/orunada/tree"
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
				for _, subspace := range subspaces{
					subspace.ComputeSubspace(x_old, x_new_update)
					subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
				}
				//os.Exit(2)
			}
		}
	}
}
