package process

import (
	"github.com/jagandecapri/vision/tree"
	"sort"
)

func ComputeDissmilarityVector(subspaces map[[2]string]tree.Subspace){
	dissimilarity_map := map[int]float64{}
	for subspace_name, subspace := range subspaces{
		center_biggest_cluster := subspace.GetBiggestCluster().GetCenter()
		outlier_clusters := subspace.GetOutliers()
		for _, cluster := range outlier_clusters{
			for _, unit := range cluster.ListOfUnits{
				for id, point := range unit.Points{
					dissimilarity_map[id] = dissimilarity_map[id] + center_biggest_cluster.Distance(&point)
				}
			}
		}
	}

	sort
}