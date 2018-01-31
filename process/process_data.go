package process

import (
	"github.com/jagandecapri/vision/tree"
	"strings"
	"github.com/jagandecapri/vision/server"
)

func GetVisualizationData(subspace tree.Subspace) (string, server.PointCluster){

	grid := subspace.Grid
	key := strings.Join(subspace.Subspace_key[:], "-")
	point_cluster := server.PointCluster{}
	for _, unit := range grid.Store{
		cluster_id := unit.Cluster_id
		tmp := [][]float64{}
		for _, point_container := range unit.Points{
			tmp = append(tmp, point_container.Vec)
		}
		point_cluster[cluster_id] = append(point_cluster[cluster_id], tmp...)
	}

	return key, point_cluster
}
