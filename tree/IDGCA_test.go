package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"math/rand"
	"math"
	"fmt"
	"sync"
	"log"
)

func TestCluster2by2Grid(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Center: Point{Vec: []float64{0.5,1.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Center: Point{Vec: []float64{1.5,0.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: Point{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res := IGDCA(&grid, min_dense_points, min_cluster_points)
	assert.True(t, res[r1].Cluster_id == res[r2].Cluster_id && res[r2].Cluster_id == res[r3].Cluster_id,
		"%v %v %v", res[r1].Cluster_id, res[r2].Cluster_id, res[r3].Cluster_id)
	assert.Equal(t, 1, len(grid.GetClusters()))
}

func TestCluster2by2GridAbsorbCluster(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: Point{Vec: []float64{0.5,1.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 1, Center: Point{Vec: []float64{1.5,0.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Cluster_id: 2, Center: Point{Vec: []float64{1.5,1.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	min_dense_points := 2

	cluster1 := Cluster{
		Cluster_id: 1,
		Num_of_points: 4,
			ListOfUnits: map[Range]*Unit{
			r2: &u2,
			r3: &u3,
		},
	}

	cluster2 := Cluster{
		Cluster_id:    2,
		Num_of_points: 2,
		ListOfUnits: map[Range]*Unit{
			r4: &u4,
		},
	}

	grid.AddUpdateCluster(cluster1)
	grid.AddUpdateCluster(cluster2)

	res, cluster, cluster_ids_to_merge := AbsorbIntoCluster(&grid, &u1, min_dense_points)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, len(cluster_ids_to_merge), "%v", cluster_ids_to_merge)
	if cluster.Cluster_id == 1{
		assert.Contains(t, cluster_ids_to_merge, 2)
		assert.Equal(t, 1, u1.Cluster_id)
		assert.Equal(t, 9, cluster.Num_of_points)
	} else {
		assert.Contains(t, cluster_ids_to_merge, 1)
		assert.Equal(t, 2, u1.Cluster_id)
		assert.Equal(t, 7, cluster.Num_of_points)
	}
}

func TestCluster2by2GridMergeClusters1(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Cluster_id: 1, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: Point{Vec: []float64{0.5,1.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 2, Center: Point{Vec: []float64{1.5,0.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: Point{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	c1 := Cluster{
		Cluster_id: 1,
		Num_of_points: 7,
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

	cluster_ids_to_be_merged := []int{2}
	res, cluster := MergeClusters(&grid, c1, cluster_ids_to_be_merged)
	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 1, u1.Cluster_id)
	assert.Equal(t, 1, u2.Cluster_id)
	assert.Equal(t, 1, u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)

	assert.Equal(t, 9, cluster.Num_of_points)

	_, ok := grid.GetCluster(2)
	assert.False(t, ok)
}

func TestCluster2by2GridMergeClusters2(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Cluster_id: 1, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Cluster_id: 1, Center: Point{Vec: []float64{0.5,1.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u3 := Unit{Id: 3, Cluster_id: 2, Center: Point{Vec: []float64{1.5,0.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u4 := Unit{Id: 4, Center: Point{Vec: []float64{1.5,1.5}}, Range: r4}
	grid.AddUnit(&u4)

	grid.SetupGrid(interval_l)

	c1 := Cluster{
		Cluster_id: 1,
		Num_of_points: 7,
			ListOfUnits: map[Range]*Unit{
			r1: &u1,
			r2: &u2,
		},
	}

	c2 := Cluster{
		Cluster_id: 2,
		Num_of_points: 2,
			ListOfUnits: map[Range]*Unit{
			r3: &u3,
		},
	}

	grid.AddUpdateCluster(c1)
	grid.AddUpdateCluster(c2)

	cluster_ids_to_be_merged := []int{1}
	res, cluster := MergeClusters(&grid, c2, cluster_ids_to_be_merged)

	assert.Equal(t, res, SUCCESS)
	assert.Equal(t, 2, u1.Cluster_id)
	assert.Equal(t, 2, u2.Cluster_id)
	assert.Equal(t, 2, u3.Cluster_id, "%v", u3.Cluster_id)
	assert.Equal(t, UNCLASSIFIED, u4.Cluster_id)

	assert.Equal(t, 9, cluster.Num_of_points)

	_, ok := grid.GetCluster(1)
	assert.False(t, ok)
}

func TestCluster3by3Grid(t *testing.T) {
	grid := NewGrid()
	interval_l := 1.0

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r1}
	grid.AddUnit(&u1)

	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	u2 := Unit{Id: 2, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {},2: {}}, Range: r2}
	grid.AddUnit(&u2)

	r3 := Range{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}
	u3 := Unit{Id: 3, Center: Point{Vec: []float64{0.5,2.5}},
		Points: map[int]Point{1: {},2: {},3: {}}, Range: r3}
	grid.AddUnit(&u3)

	r4 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	u4 := Unit{Id: 4, Center: Point{Vec: []float64{1.5,0.5}},
		Points: map[int]Point{1: {}}, Range: r4}
	grid.AddUnit(&u4)

	r5 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u5 := Unit{Id: 5, Center: Point{Vec: []float64{1.5,1.5}}, Range: r5}
	grid.AddUnit(&u5)

	r6 := Range{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}
	u6 := Unit{Id: 6, Center: Point{Vec: []float64{1.5,2.5}}, Range: r6}
	grid.AddUnit(&u6)

	r7 := Range{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}
	u7 := Unit{Id: 7, Center: Point{Vec: []float64{2.5,0.5}}, Range: r7}
	grid.AddUnit(&u7)

	r8 := Range{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}
	u8 := Unit{Id: 8, Center: Point{Vec: []float64{2.5,1.5}},
		Points: map[int]Point{1: {},2: {}, 3: {},4: {}}, Range: r8}
	grid.AddUnit(&u8)

	r9 := Range{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}
	u9 := Unit{Id: 9, Center: Point{Vec: []float64{0.5,0.5}},
		Points: map[int]Point{1: {}}, Range: r9}
	grid.AddUnit(&u9)

	grid.SetupGrid(interval_l)

	min_dense_points := 2
	min_cluster_points := 5

	res := IGDCA(&grid, min_dense_points, min_cluster_points)
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

	assert.Equal(t, 2, len(grid.GetClusters()))

}

func TestGDA(t *testing.T){
	grid := NewGrid()
	interval_l := 1.0

	for i:= 0; i < 5; i++{
		i_float := float64(i)
		r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
		u := Unit{Id: i, Cluster_id: UNCLASSIFIED, Center: Point{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
			Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r}
		grid.AddUnit(&u)
	}
	grid.SetupGrid(interval_l)
	min_dense_points := 2
	min_cluster_points := 5

	IGDCA(&grid, min_dense_points, min_cluster_points)
	assert.Equal(t, 1, len(grid.GetClusters()))

}

func BenchmarkIGDCA(t *testing.B) {
	for i:=0; i < t.N; i++{
		grid := NewGrid()
		interval_l := 1.0

		for i:= 0; i < 5; i++{
			i_float := float64(i)
			r := Range{Low: [2]float64{i_float, 0}, High: [2]float64{i_float + 1.0, 1.0}}
			u := Unit{Id: i, Center: Point{Vec: []float64{(( i_float + (i_float + 1.0))/2.0),0.5}},
				Points: map[int]Point{1: {},2: {},3: {},4: {},5: {}}, Range: r}
			grid.AddUnit(&u)
		}
		grid.SetupGrid(interval_l)
		min_dense_points := 2
		min_cluster_points := 5
		IGDCA(&grid, min_dense_points, min_cluster_points)
		assert.Equal(t, 1, len(grid.GetClusters()))
	}
}

func setupGrids(j int) (grids []*Grid, num_grids, min, max, min_dense_points, min_cluster_points int) {
	min = int(math.Pow(10.0, float64(j)))
	max = int(math.Pow(10.0, float64(j+1)))
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	min_dense_points = 10
	min_cluster_points = 100
	num_grids = 100
	grids = []*Grid{}

	for i := 0; i < num_grids; i++{
		grid := NewGrid()
		ranges := RangeBuilder(min_interval, max_interval, interval_length)
		UnitsBuilder(ranges, dim)
		units := UnitsBuilder(ranges, dim)
		rand1 := rand.New(rand.NewSource(int64(i + 1)))
		for rg, unit := range units{
			num := rand1.Intn((max - min) + 1) + min
			for i := 0; i < num; i++{
				unit.Points[i] = Point{}
			}
			units[rg] = unit
		}

		for _, unit := range units{
			grid.AddUnit(unit)
		}
		grid.SetupGrid(interval_length)
		grids = append(grids, &grid)
	}
	return
}

func BenchmarkIGDCASequential(b *testing.B){
	for j := 0; j <= 3; j++{
		b.StopTimer()
		grids, num_grids, min, max, min_dense_points, min_cluster_points := setupGrids(j)
		b.StartTimer()

		b.Run(fmt.Sprintf("Sequential/min-max=%d-%d", min, max), func(b *testing.B){
			for n := 0; n < b.N; n++{
				for m := 0; m < num_grids; m++{
					grid := grids[m]
					IGDCA(grid, min_dense_points, min_cluster_points)
				}
			}
		})
	}
}

func BenchmarkIGDCAConcurrent(b *testing.B){
	for j := 0; j <= 3; j++{
		b.StopTimer()
		grids, num_grids, min, max, min_dense_points, min_cluster_points := setupGrids(j)
		b.StartTimer()

		b.Run(fmt.Sprintf("Concurrent/min-max=%d-%d", min, max), func(b *testing.B){
			for n := 0; n < b.N; n++{
				wg := sync.WaitGroup{}
				wg.Add(num_grids)

				for m := 0; m < num_grids; m++{
					grid := grids[m]
					go func(){
						IGDCA(grid, min_dense_points, min_cluster_points)
						wg.Done()
						return
					}()
				}

				wg.Wait()
			}
		})
	}
}

func TestCustomBenchmarkIDGCASequential(t *testing.T){
	for j := 0; j <= 3; j++{
		iteration := 10000
		times_elapsed := make([]int64, iteration)

		var min, max int

		for n := 0; n < iteration; n++{
			grids, num_grids, min_tmp, max_tmp, min_dense_points, min_cluster_points := setupGrids(j)
			min = min_tmp
			max = max_tmp

			func(){
				defer func(start time.Time, name string){
					elapsed := time.Since(start)
					times_elapsed = append(times_elapsed, int64(elapsed))
				}(time.Now(), "")
				for m := 0; m < num_grids; m++{
					grid := grids[m]
					IGDCA(grid, min_dense_points, min_cluster_points)
				}
				return
			}()
		}

		var total int64
		for _, time := range times_elapsed{
			total += time
		}

		average := total/int64(iteration)
		log.Printf("Sequential/min-max=%d-%d Iteration/%d Average_Time(ns)/%d", min, max, iteration, average)
	}
}

func TestCustomBenchmarkIDGCAConcurrent(t *testing.T){
	for j := 0; j <= 3; j++{
		iteration := 10000
		times_elapsed := make([]int64, iteration)

		var min, max int

		for n := 0; n < iteration; n++{
			grids, num_grids, min_tmp, max_tmp, min_dense_points, min_cluster_points := setupGrids(j)
			min = min_tmp
			max = max_tmp

			func(){
				defer func(start time.Time, name string){
					elapsed := time.Since(start)
					times_elapsed = append(times_elapsed, int64(elapsed))
				}(time.Now(), "")

				wg := sync.WaitGroup{}
				wg.Add(num_grids)

				for m := 0; m < num_grids; m++{
					grid := grids[m]
					go func(){
						IGDCA(grid, min_dense_points, min_cluster_points)
						wg.Done()
						return
					}()
				}

				wg.Wait()
				return
			}()
		}

		var total int64
		for _, time := range times_elapsed{
			total += time
		}

		average := total/int64(iteration)
		log.Printf("Concurrent/min-max=%d-%d Iteration/%d Average_Time(ns)/%d", min, max, iteration, average)
	}
}