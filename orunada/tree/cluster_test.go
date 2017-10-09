package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCluster2by2Grid(t *testing.T) {
	units := Units{Store: make(map[Range]*Unit)}
	interval_l := 1.0

	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	units.AddUnit(&u1, r1)

	u2 := Unit{Id: 2, Center: PointContainer{Vec: []float64{0.5,1.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	units.AddUnit(&u2, r2)

	u3 := Unit{Id: 3, Center: PointContainer{Vec: []float64{1.5,0.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	units.AddUnit(&u3, r3)

	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}}
	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.AddUnit(&u4, r4)

	units.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res, cluster_map := GDA(units.Store, min_dense_points, min_cluster_points)
	assert.True(t, res[r1].Cluster_id == res[r2].Cluster_id,
		"%v %v %v", res[r1].Cluster_id, res[r2].Cluster_id)
	assert.Equal(t, 1, len(cluster_map["non_outlier_clusters"]), "%v", cluster_map["non_outlier_clusters"])
	assert.Equal(t, 0, len(cluster_map["outlier_clusters"]), "%v", cluster_map["outlier_clusters"])
}

func TestCluster3by3Grid(t *testing.T) {
	units := Units{Store: make(map[Range]*Unit)}
	interval_l := 1.0

	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	units.AddUnit(&u1, r1)

	u2 := Unit{Id: 2, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	units.AddUnit(&u2, r2)

	u3 := Unit{Id: 3, Center: PointContainer{Vec: []float64{0.5,2.5}},
		points: map[int]PointContainer{1: {},2: {},3: {}}}
	r3 := Range{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}
	units.AddUnit(&u3, r3)

	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,0.5}},
		points: map[int]PointContainer{1: {}}}
	r4 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	units.AddUnit(&u4, r4)

	u5 := Unit{Id: 5, Center: PointContainer{Vec: []float64{1.5,1.5}}}
	r5 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.AddUnit(&u5, r5)

	u6 := Unit{Id: 6, Center: PointContainer{Vec: []float64{1.5,2.5}}}
	r6 := Range{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}
	units.AddUnit(&u6, r6)

	u7 := Unit{Id: 7, Center: PointContainer{Vec: []float64{2.5,0.5}}}
	r7 := Range{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}
	units.AddUnit(&u7, r7)

	u8 := Unit{Id: 8, Center: PointContainer{Vec: []float64{2.5,1.5}},
		points: map[int]PointContainer{1: {},2: {}, 3: {},4: {}}}
	r8 := Range{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}
	units.AddUnit(&u8, r8)

	u9 := Unit{Id: 9, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {}}}
	r9 := Range{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}
	units.AddUnit(&u9, r9)

	units.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res, cluster_map := GDA(units.Store, min_dense_points, min_cluster_points)
	assert.True(t, res[r1].Cluster_id == res[r2].Cluster_id && res[r2].Cluster_id == res[r3].Cluster_id,
	"%v %v %v", res[r1].Cluster_id, res[r2].Cluster_id, res[r3].Cluster_id)
	for _, unit := range res{
		switch unit.Id{
		case 1:
		case 2:
		case 3:
		case 8:
		default:
			assert.Equal(t, 0, unit.Cluster_id, "%v", unit.Id)
		}
	}
	assert.Equal(t, 1, len(cluster_map["non_outlier_clusters"]), "%v", cluster_map["non_outlier_clusters"])
	assert.Equal(t, 1, len(cluster_map["outlier_clusters"]), "%v", cluster_map["outlier_clusters"])
}
