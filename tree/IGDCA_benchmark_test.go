package tree

import (
	"testing"
	"sync"
	"math/rand"
	"fmt"
	"math"
)

func setupGrids(min_points int, max_points int, num_grids int) (grids []*Grid, min_dense_points, min_cluster_points int) {
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	dim := 2
	min_dense_points = 10
	min_cluster_points = 100
	grids = []*Grid{}

	for i := 0; i < num_grids; i++{
		grid := NewGrid()
		ranges := RangeBuilder(min_interval, max_interval, interval_length)
		UnitsBuilder(ranges, dim)
		units := UnitsBuilder(ranges, dim)
		rand1 := rand.New(rand.NewSource(int64(i + 1)))
		for rg, unit := range units{
			num := rand1.Intn((max_points - min_points) + 1) + min_points
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

func BenchmarkIGDCAVaryPoints(b *testing.B){
	tests := []struct{
		name string
		fun func(*testing.B, []*Grid, int, int)
	}{
		{"Sequential", Sequential},
		{"Concurrent", Concurrent},
		{"Worker", Worker},
	}

	for _, test := range tests{
		for j := 0; j <= 3; j++{
			num_grids := 8
			min := int(math.Pow(10.0, float64(j)))
			max := int(math.Pow(10.0, float64(j+1)))
			b.Run(fmt.Sprintf("%v/min-max=%d-%d/num_grids=%d", test.name, min, max, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(min, max, 8)
				b.StartTimer()
				for n := 0; n < b.N; n++{
					test.fun(b, grids, min_dense_points, min_cluster_points)
				}
			})
		}
	}
}

func BenchmarkIGDCAVarySubspaces(b *testing.B){
	tests := []struct{
		name string
		fun func(*testing.B, []*Grid, int, int)
	}{
		{"Sequential", Sequential},
		{"Concurrent", Concurrent},
		{"Worker", Worker},
	}

	for _, test := range tests{
		for j := 0; j < 7; j++{
			num_grids := int(math.Pow(2.0, float64(j)))
			min := 10
			max := 100
			b.Run(fmt.Sprintf("%v/min-max=%d-%d/num_grids=%d", test.name, min, max, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(10, 100,  num_grids)
				b.StartTimer()
				for n := 0; n < b.N; n++{
					test.fun(b, grids, min_dense_points, min_cluster_points)
				}
			})
		}
	}
}

func Sequential(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		IGDCA(grid, min_dense_points, min_cluster_points)
	}
}

func Concurrent(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	wg := sync.WaitGroup{}
	wg.Add(len(grids))

	for m := 0; m < len(grids); m++{
		grid := grids[m]
		go func(){
			IGDCA(grid, min_dense_points, min_cluster_points)
			wg.Done()
			return
		}()
	}

	wg.Wait()
}

func Worker(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	b.StopTimer()
	jobs := make(chan *Grid, len(grids))
	results := make(chan struct{}, len(grids))

	for w := 1; w <= 10; w++ {
		go worker(w, min_dense_points, min_cluster_points, jobs, results)
	}
	b.StartTimer()

	for m := 0; m < len(grids); m++{
		grid := grids[m]
		jobs <- grid
	}
	close(jobs)

	for m := 0; m < len(grids); m++{
		<-results
	}
}

func worker(id int, min_dense_points int, min_cluster_points int, jobs <-chan *Grid, results chan<- struct{}) {
	for j := range jobs {
		IGDCA(j, min_dense_points, min_cluster_points)
		results <- struct{}{}
	}
}
