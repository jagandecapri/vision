package process

import (
	"testing"
	"github.com/jagandecapri/vision/tree"
	"github.com/stretchr/testify/assert"
	"github.com/jagandecapri/vision/server"
)

func TestMarshalData(t *testing.T){
	point := tree.Point{Id: 1, Vec_map: map[string]float64{
		"col_1": 0.12345,
		"col_2": 0.56789,

	}}

	point_container := tree.PointContainer{Unit_id: 5,
		Vec: []float64{0.12345, 0.56789},
		Point: point}

	center_point_container := tree.PointContainer{Unit_id: 5, Vec: []float64{0.12345, 0.56789}}

	range_1 := tree.Range{Low: [2]float64{0.1, 0.5},
		High: [2]float64{0.2, 0.6}}

	unit_1 := tree.Unit{Id: 5,
		Cluster_id: 3,
		Dimension: 2,
		Center: center_point_container,
		Points: map[int]tree.PointContainer{1: point_container},
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