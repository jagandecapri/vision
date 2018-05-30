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
		benchmark BenchmarkInterface
	}{
		{"Sequential", &SequentialExecutor{}},
		{"Concurrent", &ConcurrentExecutor{}},
		{"WorkerPoolExecutor", &WorkerPoolExecutor{}},
		{"WorkerPoolUniqueResChannelExecutor", & WorkerPoolUniqueResChannelExecutor{}},
	}

	for _, test := range tests{
		for j := 0; j <= 3; j++{
			num_grids := 8
			min := int(math.Pow(10.0, float64(j)))
			max := int(math.Pow(10.0, float64(j+1)))
			b.Run(fmt.Sprintf("%v/min-max=%d-%d/num_grids=%d", test.name, min, max, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(min, max, 8)
				test.benchmark.Init(grids, min_dense_points, min_cluster_points)
				b.StartTimer()
				for n := 0; n < b.N; n++{
					test.benchmark.Run(b, grids, min_dense_points, min_cluster_points)
				}
				b.StopTimer()
				test.benchmark.Clean()
				b.StartTimer()
			})
		}
	}
}

func BenchmarkIGDCAVarySubspaces(b *testing.B){
	tests := []struct{
		name string
		benchmark BenchmarkInterface
	}{
		{"Sequential", &SequentialExecutor{}},
		{"Concurrent", &ConcurrentExecutor{}},
		{"WorkerPoolExecutor", &WorkerPoolExecutor{}},
		{"WorkerPoolUniqueResChannelExecutor", & WorkerPoolUniqueResChannelExecutor{}},
	}

	for _, test := range tests{
		for j := 0; j < 7; j++{
			num_grids := int(math.Pow(2.0, float64(j)))
			min := 10
			max := 100
			b.Run(fmt.Sprintf("%v/min-max=%d-%d/num_grids=%d", test.name, min, max, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(10, 100,  num_grids)
				test.benchmark.Init(grids, min_dense_points, min_cluster_points)
				b.StartTimer()
				for n := 0; n < b.N; n++{
					test.benchmark.Run(b, grids, min_dense_points, min_cluster_points)
				}
				b.StopTimer()
				test.benchmark.Clean()
				b.StartTimer()
			})
		}
	}
}

type BenchmarkInterface interface{
	Init([]*Grid, int, int)
	Run(*testing.B, []*Grid, int, int)
	Clean()
}

type SequentialExecutor struct{}

func (s *SequentialExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int){}

func (s *SequentialExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		IGDCA(grid, min_dense_points, min_cluster_points)
	}
}

func (s *SequentialExecutor) Clean(){}

type ConcurrentExecutor struct{}

func (c *ConcurrentExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int){}

func (c *ConcurrentExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
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

func (c *ConcurrentExecutor) Clean(){}

type WorkerPoolExecutor struct{
	jobs chan *Grid
	results chan struct{}
}

func (w *WorkerPoolExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int){
	w.jobs = make(chan *Grid, len(grids))
	w.results = make(chan struct{}, len(grids))
	num_workers := 4
	for i := 0; i < num_workers; i++ {
		go Worker(i, min_dense_points, min_cluster_points, w.jobs, w.results)
	}
}

func (w *WorkerPoolExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		w.jobs <- grid
	}

	for m := 0; m < len(grids); m++{
		<-w.results
	}
}

func Worker(id int, min_dense_points int, min_cluster_points int, jobs <-chan *Grid, results chan<- struct{}) {
	for{
		select{
			case j, ok := <-jobs:
				if ok{
					IGDCA(j, min_dense_points, min_cluster_points)
					results <- struct{}{}
				} else {
					return
				}
			default:
		}
	}
}


func (w *WorkerPoolExecutor) Clean(){
	close(w.jobs)
	close(w.results)
}

type WorkerPoolUniqueResChannelExecutor struct{
	jobs chan *Grid
	done chan struct{}
	final_out <-chan struct{}
}

func (w *WorkerPoolUniqueResChannelExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int){
	w.jobs = make(chan *Grid, len(grids))
	w.done = make(chan struct{})
	num_workers := 4

	var cs []chan struct{}
	for i := 0; i < num_workers; i++ {
		c := WorkerUnique(i, min_dense_points, min_cluster_points, w.done, w.jobs)
		cs = append(cs, c)
	}

	w.final_out = Receiver(w.done, cs...)
}

func (w *WorkerPoolUniqueResChannelExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		w.jobs <- grid
	}

	for m := 0; m < len(grids); m++{
		<-w.final_out
	}
}

func (w *WorkerPoolUniqueResChannelExecutor) Clean(){
	close(w.done)
}


func WorkerUnique(id int, min_dense_points int, min_cluster_points int, done <-chan struct{}, jobs <-chan *Grid) chan struct{}{
	out := make(chan struct{})
	go func(){
		defer close(out)
		for{
			select{
			case j := <-jobs:
					IGDCA(j, min_dense_points, min_cluster_points)
					out <- struct{}{}
			case <-done:
				return
			}
		}
	}()
	return out
}

func Receiver(done <-chan struct{}, cs ...chan struct{}) <-chan struct{}{
	wg := sync.WaitGroup{}
	out := make(chan struct{})

	output := func(c <-chan struct{}){
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}