package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func CountClusterTypes(cluster_list map[int]Cluster) (int, int){
	outlier_cluster_count := 0
	non_outlier_cluster_count := 0
	for _, cluster := range cluster_list{
		if cluster.Cluster_type == OUTLIER_CLUSTER{
			outlier_cluster_count++
		} else if cluster.Cluster_type == NON_OUTLIER_CLUSTER{
			non_outlier_cluster_count++
		}
	}

	return outlier_cluster_count, non_outlier_cluster_count
}

func TestCluster2by2Grid(t *testing.T) {
	units := NewUnits()
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

	res, cluster_map := GDA(units, min_dense_points, min_cluster_points)
	assert.True(t, res[r1].Cluster_id == res[r2].Cluster_id && res[r2].Cluster_id == res[r3].Cluster_id,
		"%v %v %v", res[r1].Cluster_id, res[r2].Cluster_id, res[r3].Cluster_id)
	outlier_cluster_count, non_outlier_cluster_count := CountClusterTypes(cluster_map)
	assert.Equal(t, 0, outlier_cluster_count)
	assert.Equal(t, 1, non_outlier_cluster_count)
}

func TestCluster2by2GridAbsorbCluster(t *testing.T) {
	units := NewUnits()
	interval_l := 1.0

	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	units.AddUnit(&u1, r1)

	u2 := Unit{Id: 2, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,1.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	units.AddUnit(&u2, r2)

	u3 := Unit{Id: 3, Cluster_id: 1, Center: PointContainer{Vec: []float64{1.5,0.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	units.AddUnit(&u3, r3)

	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}}
	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.AddUnit(&u4, r4)

	units.SetupGrid(interval_l)

	min_dense_points := 2

	res, cluster_ids := AbsorbIntoCluster(&u1, min_dense_points)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, len(cluster_ids), "%v", cluster_ids)
	assert.Contains(t, cluster_ids, 1)
	assert.Equal(t, 1, u1.Cluster_id)
}

func TestCluster2by2GridMergeClusters(t *testing.T) {
	units := NewUnits()
	interval_l := 1.0

	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	units.AddUnit(&u1, r1)

	u2 := Unit{Id: 2, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,1.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	units.AddUnit(&u2, r2)

	u3 := Unit{Id: 3, Cluster_id: 2, Center: PointContainer{Vec: []float64{1.5,0.5}},
		points: map[int]PointContainer{1: {},2: {}}}
	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	units.AddUnit(&u3, r3)

	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}}
	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.AddUnit(&u4, r4)

	units.SetupGrid(interval_l)

	units.Cluster_map = map[int]Cluster{
		1: {
			Cluster_id: 1,
			ListOfUnits: map[Range]*Unit{
				r1: &u1,
				r2: &u2,
			},
		},
		2: {
			Cluster_id: 2,
			ListOfUnits: map[Range]*Unit{
				r3: &u3,
			},
		},
	}

	var cluster_ids []int
	var res int
	cluster_ids = []int{1,2}
	res = MergeClusters(units.Cluster_map,cluster_ids)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, u1.Cluster_id)
	assert.Equal(t, 1, u2.Cluster_id)
	assert.Equal(t, 1, u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)

	cluster_ids = []int{2,1}
	res = MergeClusters(units.Cluster_map,cluster_ids)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 2, u1.Cluster_id)
	assert.Equal(t, 2, u2.Cluster_id)
	assert.Equal(t, 2, u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)
}

func TestCluster3by3Grid(t *testing.T) {
	units := NewUnits()
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

	res, cluster_map := GDA(units, min_dense_points, min_cluster_points)
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
	outlier_cluster_count, non_outlier_cluster_count := CountClusterTypes(cluster_map)
	assert.Equal(t, 1, outlier_cluster_count)
	assert.Equal(t, 1, non_outlier_cluster_count)
}

func TestGDA(t *testing.T){
	units := NewUnits()
	interval_l := 1.0

	for i:= 0; i < 5; i++{
		i_float := float64(i)
		u := Unit{Id: i, Cluster_id: UNCLASSIFIED, Center: PointContainer{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
			points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
		r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
		units.AddUnit(&u, r)
	}
	units.SetupGrid(interval_l)
	min_dense_points := 2
	min_cluster_points := 5

	_, cluster_map := GDA(units, min_dense_points, min_cluster_points)
	outlier_cluster_count, non_outlier_cluster_count := CountClusterTypes(cluster_map)
	assert.Equal(t, 0, outlier_cluster_count)
	assert.Equal(t, 1, non_outlier_cluster_count)
}

func BenchmarkGDA(t *testing.B) {
	for i:=0; i < t.N; i++{
		units := Units{Store: make(map[Range]*Unit)}
		interval_l := 1.0

		for i:= 0; i < 5; i++{
			i_float := float64(i)
			u := Unit{Id: i, Center: PointContainer{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
				points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}}
			r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
			units.AddUnit(&u, r)
		}
		units.SetupGrid(interval_l)
		min_dense_points := 2
		min_cluster_points := 5
		_, cluster_map := GDA(units, min_dense_points, min_cluster_points)
		outlier_cluster_count, non_outlier_cluster_count := CountClusterTypes(cluster_map)
		assert.Equal(t, 0, outlier_cluster_count)
		assert.Equal(t, 1, non_outlier_cluster_count)
	}
}
