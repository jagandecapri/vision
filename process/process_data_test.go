package process

import (
	"testing"
	"github.com/jagandecapri/vision/tree"
	"github.com/stretchr/testify/assert"
	"github.com/jagandecapri/vision/server"
	"encoding/json"
)

func TestProcessDataForVisualization(t *testing.T) {

	r1 := tree.Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := tree.Unit{Id: 1, Center: tree.Point{Vec: []float64{0.5,0.5}},
		Points: map[int]tree.Point{1: {Vec: []float64{0.5, 0.5}}}, Range: r1}


	r2 := tree.Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u2 := tree.Unit{Id: 3, Center: tree.Point{Vec: []float64{1.5,1.5}},
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

	subspace_key1 := [2]string{"third", "fourth"}

	subspace1 := tree.Subspace{
		Grid: &grid,
		Subspace_key: subspace_key1,
	}

	subspaces := []tree.Subspace{subspace, subspace1}

	expected_struct := server.HttpData1{
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "first-second"},
			Points: []server.Point{
				server.Point{
					Point_data: server.Point_data{X: 0.5,Y: 0.5},
					Point_metadata: server.Point_metadata{Color: "#ABC"},
				},
				server.Point{
					Point_data: server.Point_data{X: 1.5,Y: 1.5},
					Point_metadata: server.Point_metadata{Color: "#ABC"},
				},
			},
		},
		server.Graph{
			Graph_metadata: server.Graph_metadata{ID: "third-fourth"},
			Points: []server.Point{
				server.Point{
					Point_data: server.Point_data{X: 0.5,Y: 0.5},
					Point_metadata: server.Point_metadata{Color: "#ABC"},
				},
				server.Point{
					Point_data: server.Point_data{X: 1.5,Y: 1.5},
					Point_metadata: server.Point_metadata{Color: "#ABC"},
				},
			},
		},
	}

	res := ProcessDataForVisualization(subspaces)
	json1, _ := json.Marshal(expected_struct)
	json2, _ := json.Marshal(res)

	assert.JSONEq(t, string(json1), string(json2))
}