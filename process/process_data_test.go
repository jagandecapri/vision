package process

import (
	"testing"
	"github.com/jagandecapri/vision/tree"
	"github.com/stretchr/testify/assert"
	"github.com/jagandecapri/vision/server"
	"github.com/stretchr/testify/mock"
	"github.com/jagandecapri/vision/utils/color/mocks"
)

func TestProcessDataForVisualization(t *testing.T) {

	r1 := tree.Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := tree.Unit{Id: 1, Center: tree.Point{Vec: []float64{0.5,0.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{0.5, 0.5}}}, Range: r1}

	r2 := tree.Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u2 := tree.Unit{Id: 3, Center: tree.Point{Vec: []float64{1.5,1.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{1.5, 1.5}}}, Range: r2}

	r3 := tree.Range{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}
	u3 := tree.Unit{Id: 1, Center: tree.Point{Vec: []float64{2.5,2.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{2.5, 2.5}}}, Range: r3}

	c1 := tree.Cluster{Cluster_id: 1,
		Cluster_type: tree.OUTLIER_CLUSTER,
		Num_of_points: 10,
		ListOfUnits: map[tree.Range]*tree.Unit{r1: &u1},
	}

	c2 := tree.Cluster{Cluster_id: 3,
		Cluster_type: tree.NON_OUTLIER_CLUSTER,
		Num_of_points: 10,
		ListOfUnits: map[tree.Range]*tree.Unit{r2: &u2},
	}

	cc := tree.ClusterContainer{ListOfClusters: map[int]tree.Cluster{1: c1, 2: c2}}

	grid := tree.NewGrid()
	grid.ClusterContainer = cc
	grid.AddUnit(&u1)
	grid.AddUnit(&u2)
	grid.AddUnit(&u3)

	subspace_key := [2]string{"first", "second"}

	subspace := tree.Subspace{
		Grid: &grid,
		Subspace_key: subspace_key,
	}

	subspace_key1 := [2]string{"third", "fourth"}

	subspace1 := tree.Subspace{
		Grid: &grid,
		Subspace_key: subspace_key1,
	}

	subspaces := []tree.Subspace{subspace, subspace1}

	/**
	Expected JSON
[
  {"metadata":{
    "id": "first-second",
    "column_x": "first",
    "column_y": "second"
  },
    "points": [
      {"data": [{
        "x": "0.5",
        "y": "0.5"
      }],
        "metadata": {
          "color": "#ABC"
        }},
      {"data": [{
        "x": "1.5",
        "y": "1.5"
      }],
        "metadata": {
          "color": "#DEF"
        }},
      {"data": [{
        "x": "2.5",
        "y": "2.5"
      }],
        "metadata": {
          "color": "#GHI"
        }}
    ]},
  {"metadata":{
    "id": "third-four"
    "column_x": "third",
    "column_y": "fourth"
  },
    "points": [
      {"data": [{
        "x": "0.5",
        "y": "0.5"
      }],
        "metadata": {
          "color": "#ABC"
        }},
      {"data": [{
        "x": "1.5",
        "y": "1.5"
      }],
        "metadata": {
          "color": "#DEF"
        }},
      {"data": [{
        "x": "2.5",
        "y": "2.5"
      }],
        "metadata": {
          "color": "#GHI"
        }}
    ]}
	]**/

	expected_struct := server.HttpData{
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "first-second",
				Column_x: "first",
				Column_y: "second"},
			Points: []server.Points{{
				Point_list: []server.Point{{
					Point_data: server.Point_data{X: 0.5, Y: 0.5},
				}},
				Point_metadata: server.Point_metadata{Color: "#ABC"},
			},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 1.5, Y: 1.5},
					}},
					Point_metadata: server.Point_metadata{Color: "#DEF"},
				},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 2.5, Y: 2.5},
					}},
					Point_metadata: server.Point_metadata{Color: "#GHI"},
				},
			},
		},
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "third-fourth",
				Column_x: "third",
				Column_y: "fourth"},
			Points: []server.Points{{
				Point_list: []server.Point{{
					Point_data: server.Point_data{X: 0.5, Y: 0.5},
				}},
				Point_metadata: server.Point_metadata{Color: "#ABC"},
			},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 1.5, Y: 1.5},
					}},
					Point_metadata: server.Point_metadata{Color: "#DEF"},
				},
				{
					Point_list: []server.Point{{
						Point_data: server.Point_data{X: 2.5, Y: 2.5},
					}},
					Point_metadata: server.Point_metadata{Color: "#GHI"},
				},
			},
		},
	}

	mock_color_helper := &mocks.ColorHelperInterface{}
	mock_color_helper.On("GetRandomColors", mock.Anything).Return([]string{"#ABC", "#DEF", "#GHI"})
	res := processDataForVisualization(subspaces, mock_color_helper)
	for i := 0; i < len(res); i++{
		tmp := res[i]
		graph := expected_struct[i]
		assert.Equal(t, graph.Graph_metadata, tmp.Graph_metadata)
		for _, points := range graph.Points{
			assert.Contains(t, tmp.Points, points)
		}
	}
}