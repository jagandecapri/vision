package process

import (
	"github.com/jagandecapri/vision/tree"
	"strings"
	"github.com/jagandecapri/vision/server"
)

func ProcessUnits(unit *tree.Unit) []server.Point{
	points := unit.GetPoints()

	tmp := []server.Point{}

	for _, point := range points{
		X := point.GetValue(0)
		Y := point.GetValue(1)

		point_data := server.Point{
			Point_data: server.Point_data{
				X: X,
				Y: Y,
			},
		}

		tmp = append(tmp, point_data)
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
		clustered_units_acc := make(map[tree.Range]*tree.Unit)

		//Clustered Unit data
		for _, cluster := range clusters{
			units := cluster.GetUnits()
			tmp := server.Points{}

			for rg, unit := range units{
				clustered_units_acc[rg] = unit
				tmp1 := ProcessUnits(unit)
				tmp.Point_list = tmp1
				tmp.Point_metadata = server.Point_metadata{
					Color: "#ABC",
				}
			}

			points_data = append(points_data, tmp)
		}

		all_units := subspace.GetUnits()

		//Unclustered unit data
		tmp2 := server.Points{}
		tmp2.Point_metadata = server.Point_metadata{
			Color: "#DEF",
		}

		for rg, unit := range all_units{
			if _, ok := clustered_units_acc[rg]; !ok{
				tmp := ProcessUnits(unit)
				tmp2.Point_list = append(tmp2.Point_list, tmp...)
			}
		}

		points_data = append(points_data, tmp2)

		graph.Points = points_data
		graphs = append(graphs, graph)
	}

	return graphs
}