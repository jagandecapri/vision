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

func ProcessDataForVisualization(subspaces []tree.Subspace) server.HttpData1{
	graphs := server.HttpData1{}

	for _, subspace := range subspaces{
		id := strings.Join(subspace.Subspace_key[:], "-")

		graph := server.Graph{
			Graph_metadata: server.Graph_metadata{ID: id},
		}

		 clusters := subspace.GetClusters()

		 points_data := []server.Point{}

		 for _, cluster := range clusters{
		 	units := cluster.GetUnits()

		 	for _, unit := range units{
		 		points := unit.GetPoints()

		 		for _, point := range points{
		 			X := point.GetValue(0)
		 			Y := point.GetValue(1)

		 			point_data := server.Point{
		 				Point_data: server.Point_data{
							X: X,
							Y: Y,
						},
						Point_metadata: server.Point_metadata{
							Color: "#ABC",
						},
					}

					points_data = append(points_data, point_data)
				}
			}
		 }

		 graph.Points = points_data

		 graphs = append(graphs, graph)
	}

	return graphs
}