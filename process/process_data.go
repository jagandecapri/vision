package process

import (
	"github.com/jagandecapri/vision/tree"
	"strings"
	"github.com/jagandecapri/vision/server"
)

func GetClusteredUnitPoints(cluster tree.Cluster) server.Points{

	units := cluster.GetUnits()
	tmp := server.Points{}

	for _, unit := range units{
		points := unit.GetPoints()

		tmp1 := []server.Point{}

		for _, point := range points{
			X := point.GetValue(0)
			Y := point.GetValue(1)

			point_data := server.Point{
				Point_data: server.Point_data{
					X: X,
					Y: Y,
				},
			}

			tmp1 = append(tmp1, point_data)
		}

		tmp.Point_list = tmp1
		tmp.Point_metadata = server.Point_metadata{
			Color: "#ABC",
		}
	}

	return tmp
}

func ProcessDataForVisualization(subspaces []tree.Subspace) server.HttpData{
	graphs := server.HttpData{}

	for _, subspace := range subspaces{
		id := strings.Join(subspace.Subspace_key[:], "-")

		graph := server.Graph{
			Graph_metadata: server.Graph_metadata{ID: id,
			Column_x: subspace.Subspace_key[0],
			Column_y: subspace.Subspace_key[1],
			},
		}

		points_data := []server.Points{}

		clusters := subspace.GetClusters()
		for _, cluster := range clusters{
			tmp := GetClusteredUnitPoints(cluster)
			points_data = append(points_data, tmp)
		}

		graph.Points = points_data
		graphs = append(graphs, graph)
	}

	return graphs
}