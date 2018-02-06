package process

import (
	"testing"
	"github.com/jagandecapri/vision/tree"
	"github.com/stretchr/testify/assert"
	"github.com/jagandecapri/vision/server"
	"encoding/json"
	"fmt"
)

func TestMarshalData(t *testing.T){
	//point := tree.Point{Id: 1, Vec_map: map[string]float64{
	//	"col_1": 0.12345,
	//	"col_2": 0.56789,
	//
	//}}

	point_container := tree.Point{Unit_id: 5,
		Vec: []float64{0.12345, 0.56789},
		Id: 1, Vec_map: map[string]float64{
			"col_1": 0.12345,
			"col_2": 0.56789,

		}}

	center_point_container := tree.Point{Unit_id: 5, Vec: []float64{0.12345, 0.56789}}

	range_1 := tree.Range{Low: [2]float64{0.1, 0.5},
		High: [2]float64{0.2, 0.6}}

	unit_1 := tree.Unit{Id: 5,
		Cluster_id: 3,
		Dimension: 2,
		Center: center_point_container,
		Points: map[int]tree.Point{1: point_container},
		Center_calculated: true,
		Range: range_1,
	}

	grid := tree.Grid{Store: map[tree.Range]*tree.Unit{range_1: &unit_1}}

	subspace := tree.Subspace{Grid: &grid, Subspace_key: [2]string{"a", "b"}}

	key, point_cluster := GetVisualizationData(subspace)

	expected_point_cluster := server.PointCluster{3: {{
		0.12345, 0.56789,
	}}}
	assert.Equal(t, "a-b", key)
	assert.Equal(t, expected_point_cluster, point_cluster)

	http_data := server.HttpData{key: point_cluster}

	assert.Equal(t, server.HttpData{"a-b": expected_point_cluster}, http_data)
	//res, err := json.Marshal(grid)
	//log.Println(res, err)
}

func TestProcessDataForVisualization(t *testing.T) {

	r1 := tree.Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := tree.Unit{Id: 1, Center: tree.Point{Vec: []float64{0.5,0.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{0.5, 0.5}}}, Range: r1}


	r2 := tree.Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u2 := tree.Unit{Id: 3, Center: tree.Point{Vec: []float64{1.5,0.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{1.5, 1.5}}}, Range: r2}

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

	subspace_key := [2]string{"first", "second"}

	subspace := tree.Subspace{
		Grid: &grid,
		Subspace_key: subspace_key,
	}

	subspace_key1 := [2]string{"third", "four"}

	subspace1 := tree.Subspace{
		Grid: &grid,
		Subspace_key: subspace_key1,
	}

	subspaces := []tree.Subspace{subspace, subspace1}

	expected_json := `[
  {"metadata":{
    "id": "first-second"
  },
    "points": [
      {"data": {
        "x": "0.5",
        "y": "0.5"
      },
        "metadata": {
          "color": "#ABC"
        }},
      {"data": {
        "x": "1.5",
        "y": "1.5"
      },
        "metadata": {
          "color": "#ABC"
        }}
    ]},
    {"metadata":{
      "id": "third-four"
    },
      "points": [
        {"data": {
          "x": "0.5",
          "y": "0.5"
        },
          "metadata": {
            "color": "#ABC"
          }},
        {"data": {
          "x": "1.5",
          "y": "1.5"
        },
          "metadata": {
            "color": "#DEF"
          }}
      ]}
]`

	http_data := ProcessDataForVisualization(subspaces)
	json_byte, _ := json.Marshal(http_data)
	assert.JSONEq(t, expected_json, string(json_byte))

	fmt.Println(string([]byte(expected_json)))
	fmt.Println(expected_json, string(json_byte))
}