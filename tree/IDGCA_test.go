package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCluster2by2Grid(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Center: PointContainer{Vec: []float64{0.5,1.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Center: PointContainer{Vec: []float64{1.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res := IGDCA(grid, min_dense_points, min_cluster_points)
	assert.True(t, res[r1].Cluster_id == res[r2].Cluster_id && res[r2].Cluster_id == res[r3].Cluster_id,
		"%v %v %v", res[r1].Cluster_id, res[r2].Cluster_id, res[r3].Cluster_id)
	assert.Equal(t, 0, len(grid.GetOutliers()))
	assert.Equal(t, 1, len(grid.GetNonOutliers()))
}

func TestCluster2by2GridAbsorbCluster(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,1.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 1, Center: PointContainer{Vec: []float64{1.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	min_dense_points := 2

	cluster := Cluster{
		Cluster_id: 1,
			ListOfUnits: map[Range]*Unit{
			r2: &u2,
			r3: &u3,
		},
	}
	grid.AddUpdateCluster(cluster)

	res, cluster_ids := AbsorbIntoCluster(grid, &u1, r1, min_dense_points)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, len(cluster_ids), "%v", cluster_ids)
	assert.Contains(t, cluster_ids, 1)
	assert.Equal(t, 1, u1.Cluster_id)
	c, _ := grid.GetCluster(1)
	_, ok := c.ListOfUnits[r1]
	assert.True(t, ok)
}

func TestCluster2by2GridMergeClusters1(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,1.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 2, Center: PointContainer{Vec: []float64{1.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	c1 := Cluster{
		Cluster_id: 1,
		ListOfUnits: map[Range]*Unit{
			r1: &u1,
			r2: &u2,
		},
	}

	c2 := Cluster{
		Cluster_id: 2,
			ListOfUnits: map[Range]*Unit{
			r3: &u3,
		},
	}

	grid.AddUpdateCluster(c1)
	grid.AddUpdateCluster(c2)

	var cluster_ids []int
	var res int
	var merged_cluster_id []int
	var ok bool

	cluster_ids = []int{1,2}
	res, merged_cluster_id = MergeClusters(grid,cluster_ids)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, u1.Cluster_id)
	assert.Equal(t, 1, u2.Cluster_id)
	assert.Equal(t, 1, u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)
	assert.Contains(t, merged_cluster_id, 1)
	_, ok = grid.GetCluster(2)
	assert.False(t, ok)
}

func TestCluster2by2GridMergeClusters2(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: PointContainer{Vec: []float64{0.5,1.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 2, Center: PointContainer{Vec: []float64{1.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	c1 := Cluster{
		Cluster_id: 1,
			ListOfUnits: map[Range]*Unit{
			r1: &u1,
			r2: &u2,
		},
	}

	c2 := Cluster{
		Cluster_id: 2,
			ListOfUnits: map[Range]*Unit{
			r3: &u3,
		},
	}

	grid.AddUpdateCluster(c1)
	grid.AddUpdateCluster(c2)

	var cluster_ids []int
	var res int
	var merged_cluster_id []int
	var ok bool

	cluster_ids = []int{2,1}
	res, merged_cluster_id = MergeClusters(grid,cluster_ids)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 2, u1.Cluster_id)
	assert.Equal(t, 2, u2.Cluster_id)
	assert.Equal(t, 2, u3.Cluster_id, "%v", u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)
	assert.Contains(t, merged_cluster_id, 2)
	_, ok = grid.GetCluster(1)
	assert.False(t, ok)
}

func TestCluster3by3Grid(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}
	u3 := Unit{Id: 3, Center: PointContainer{Vec: []float64{0.5,2.5}},
		Points: map[int]PointContainer{1: {},2: {},3: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u4 := Unit{Id: 4, Center: PointContainer{Vec: []float64{1.5,0.5}},
		Points: map[int]PointContainer{1: {}}, Range: r4}
	grid.AddUnit(&u4)

	r5 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u5 := Unit{Id: 5, Center: PointContainer{Vec: []float64{1.5,1.5}}, Range: r5}
	grid.AddUnit(&u5)

	r6 := Range{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}
	u6 := Unit{Id: 6, Center: PointContainer{Vec: []float64{1.5,2.5}}, Range: r6}
	grid.AddUnit(&u6)

	r7 := Range{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}
	u7 := Unit{Id: 7, Center: PointContainer{Vec: []float64{2.5,0.5}}, Range: r7}
	grid.AddUnit(&u7)

	r8 := Range{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}
	u8 := Unit{Id: 8, Center: PointContainer{Vec: []float64{2.5,1.5}},
		Points: map[int]PointContainer{1: {},2: {}, 3: {},4: {}}, Range: r8}
	grid.AddUnit(&u8)

	r9 := Range{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}
	u9 := Unit{Id: 9, Center: PointContainer{Vec: []float64{0.5,0.5}},
		Points: map[int]PointContainer{1: {}}, Range: r9}
	grid.AddUnit(&u9)

	grid.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res := IGDCA(grid, min_dense_points, min_cluster_points)
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

	assert.Equal(t, 1, len(grid.GetOutliers()))
	assert.Equal(t, 1, len(grid.GetNonOutliers()))
}

func TestGDA(t *testing.T){
	grid := NewGrid()
	interval_l := 1.0

	for i:= 0; i < 5; i++{
		i_float := float64(i)
		r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
		u := Unit{Id: i, Cluster_id: UNCLASSIFIED, Center: PointContainer{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
			Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r}
		grid.AddUnit(&u)
	}
	grid.SetupGrid(interval_l)
	min_dense_points := 2
	min_cluster_points := 5

	IGDCA(grid, min_dense_points, min_cluster_points)
	assert.Equal(t, 0, len(grid.GetOutliers()))
	assert.Equal(t, 1, len(grid.GetNonOutliers()))
}

func BenchmarkGDA(t *testing.B) {
	for i:=0; i < t.N; i++{
		grid := NewGrid()
		interval_l := 1.0

		for i:= 0; i < 5; i++{
			i_float := float64(i)
			r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
			u := Unit{Id: i, Center: PointContainer{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
				Points: map[int]PointContainer{1: {},2: {},3: {},4: {},5: {}}, Range: r}
			grid.AddUnit(&u)
		}
		grid.SetupGrid(interval_l)
		min_dense_points := 2
		min_cluster_points := 5
		IGDCA(grid, min_dense_points, min_cluster_points)
		assert.Equal(t, 0, len(grid.GetOutliers()))
		assert.Equal(t, 1, len(grid.GetNonOutliers()))
	}
}
