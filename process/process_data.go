package process

import (
	"github.com/jagandecapri/vision/tree"
	"strings"
	"github.com/jagandecapri/vision/server"
	"github.com/jagandecapri/vision/utils/color"
)

func processUnit(unit *tree.Unit) []server.Point{
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

func processDataForVisualization(subspaces map[[2]string]tree.Subspace, color_helper color.ColorHelperInterface) server.HttpData{
	graphs := server.HttpData{}

	for _, subspace := range subspaces{
		id := strings.Join(subspace.Subspace_key[:], "-")

		graph := server.Graph{
			Graph_metadata: server.Graph_metadata{ID: id,
			Column_x: subspace.Subspace_key[0],
			Column_y: subspace.Subspace_key[1],
			},
		}

		points_data := []server.PointsContainer{}

		clusters := subspace.GetClusters()
		clustered_units_acc := make(map[tree.Range]*tree.Unit)

		colors := color_helper.GetRandomColors(len(clusters) + 1) // 1 is for unclustered unit

		//Clustered Unit data
		i := 0
		for _, cluster := range clusters {
			color := colors[i]
			i++
			units := cluster.GetUnits()
			tmp := server.PointsContainer{}
			tmp.Points_metadata = server.Points_metadata{
				Color: color,
			}

			for rg, unit := range units{
				clustered_units_acc[rg] = unit
				tmp1 := processUnit(unit)
				tmp.Point_list = tmp1
			}

			points_data = append(points_data, tmp)
		}

		all_units := subspace.GetUnits()

		//Unclustered unit data
		color := colors[len(colors) - 1]
		tmp2 := server.PointsContainer{}
		tmp2.Points_metadata = server.Points_metadata{
			Color: color,
		}

		for rg, unit := range all_units{
			if _, ok := clustered_units_acc[rg]; !ok{
				tmp := processUnit(unit)
				tmp2.Point_list = append(tmp2.Point_list, tmp...)
			}
		}

		points_data = append(points_data, tmp2)

		graph.PointsContainer = points_data
		graphs = append(graphs, graph)
	}

	return graphs
}

var color_helper color.ColorHelperInterface

func ProcessDataForVisualization(subspaces map[[2]string]tree.Subspace) server.HttpData{
	if color_helper == nil{
		color_helper = &color.ColorHelper{}
	}
	return processDataForVisualization(subspaces, color_helper)
}