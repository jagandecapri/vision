package process

import (
	"github.com/jagandecapri/vision/tree"
)

type DissimilarityVector map[tree.PointKey]float64

func ComputeDissmilarityVector(subspace tree.Subspace) DissimilarityVector{
	dissimilarity_map := map[tree.PointKey]float64{}

	outliers := subspace.GetOutliers()
	cluster := subspace.GetBiggestCluster()
	center_biggest_cluster := cluster.GetCenter()

	for _, outlier := range outliers{
		distance := center_biggest_cluster.Distance(&outlier)
		dissimilarity_map[outlier.Key] = distance
	}

	return dissimilarity_map
}