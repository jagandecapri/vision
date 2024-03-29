package tree

import (
	"testing"
	"sync"
	"math/rand"
	"fmt"
	"math"
	"runtime"
	//"sync/atomic"
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
		{"WorkerPoolSharedResChannelExecutor", &WorkerPoolSharedResExecutor{}},
		{"WorkerPoolUniqueResChannelExecutor", & WorkerPoolUniqueResChannelExecutor{}},
	}

	num_cpu := runtime.NumCPU()

	for _, test := range tests{
		for j := 0; j < 4; j++{
			num_grids := 8
			min_points := int(math.Pow(10.0, float64(j)))
			max_points := int(math.Pow(10.0, float64(j+1)))
			b.Run(fmt.Sprintf("%v/min_points-max_points=%d-%d/num_grids=%d", test.name, min_points, max_points, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(min_points, max_points, 8)
				test.benchmark.Init(grids, min_dense_points, min_cluster_points, num_cpu)
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
		{"WorkerPoolSharedResChannelExecutor", &WorkerPoolSharedResExecutor{}},
		{"WorkerPoolUniqueResChannelExecutor", & WorkerPoolUniqueResChannelExecutor{}},
	}

	num_cpu := runtime.NumCPU()

	for _, test := range tests{
		for j := 0; j < 11; j++{
			num_grids := int(math.Pow(2.0, float64(j)))
			min_points := 10
			max_points := 100
			b.Run(fmt.Sprintf("%v/min_points-max_points=%d-%d/num_grids=%d", test.name, min_points, max_points, num_grids), func(b *testing.B){
				b.StopTimer()
				grids, min_dense_points, min_cluster_points := setupGrids(min_points, max_points,  num_grids)
				test.benchmark.Init(grids, min_dense_points, min_cluster_points, num_cpu)
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
	Init([]*Grid, int, int, int)
	Run(*testing.B, []*Grid, int, int)
	Clean()
}

type SequentialExecutor struct{}

func (s *SequentialExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int, num_cpu int){}

func (s *SequentialExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		IGDCA(grid, min_dense_points, min_cluster_points)
	}
}

func (s *SequentialExecutor) Clean(){}

type ConcurrentExecutor struct{
	num_cpu int
}

func (c *ConcurrentExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int, num_cpu int){
	c.num_cpu = num_cpu
}

func (c *ConcurrentExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	mod := len(grids)%c.num_cpu
	iterations := (len(grids) - (mod))/c.num_cpu

	if mod > 0{
		iterations++
	}

	for i := 0; i < iterations; i++{
		wg := sync.WaitGroup{}

		min := i * c.num_cpu
		max := (i + c.num_cpu) * c.num_cpu

		if max > len(grids){
			max = len(grids)
		}

		sub_set := grids[min:max]

		wg.Add(len(sub_set))

		for m := 0; m < len(sub_set); m++{
			grid := sub_set[m]
			go func(){
				IGDCA(grid, min_dense_points, min_cluster_points)
				wg.Done()
				return
			}()
		}

		wg.Wait()
	}
}

func (c *ConcurrentExecutor) Clean(){}

type WorkerPoolSharedResExecutor struct{
	jobs chan *Grid
	results chan struct{}
}

func (w *WorkerPoolSharedResExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int, num_cpu int){
	w.jobs = make(chan *Grid, len(grids))
	w.results = make(chan struct{}, len(grids))
	num_workers := num_cpu
	for i := 0; i < num_workers; i++ {
		go Worker(i, min_dense_points, min_cluster_points, w.jobs, w.results)
	}
}

func (w *WorkerPoolSharedResExecutor) Run(b *testing.B, grids []*Grid, min_dense_points int, min_cluster_points int){
	for m := 0; m < len(grids); m++{
		grid := grids[m]
		w.jobs <- grid
	}

	for m := 0; m < len(grids); m++{
		<-w.results
	}

	//b.Log("defaults", atomic.LoadUint64(&defaults))
	//b.Log("proceeds", atomic.LoadUint64(&proceeds))
}

//var defaults uint64
//var proceeds uint64

func Worker(id int, min_dense_points int, min_cluster_points int, jobs <-chan *Grid, results chan<- struct{}) {
	for{
		select{
			case j, ok := <-jobs:
				//atomic.AddUint64(&proceeds, 1)
				if ok{
					IGDCA(j, min_dense_points, min_cluster_points)
					results <- struct{}{}
				} else {
					return
				}
			//default:
				//atomic.AddUint64(&defaults, 1)
		}
	}
}


func (w *WorkerPoolSharedResExecutor) Clean(){
	close(w.jobs)
	close(w.results)
}

type WorkerPoolUniqueResChannelExecutor struct{
	jobs chan *Grid
	done chan struct{}
	final_out <-chan struct{}
}

func (w *WorkerPoolUniqueResChannelExecutor) Init(grids []*Grid, min_dense_points int, min_cluster_points int, num_cpu int){
	w.jobs = make(chan *Grid, len(grids))
	w.done = make(chan struct{})
	num_workers := num_cpu

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