package process

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/tree"
	"sync"
	"fmt"
	"runtime"
	"time"
	"github.com/jagandecapri/vision/utils"
	"github.com/jagandecapri/vision/anomalies"
	"github.com/jagandecapri/vision/server"
)

type AccumulatorChannel chan preprocess.MicroSlot

type AccumulatorChannels struct{
	AggSrc AccumulatorChannel
	AggDst AccumulatorChannel
	AggSrcDst AccumulatorChannel
}

func UpdateFeatureSpaceBuilder(subspace_channel_container anomalies.SubspaceChannelsContainer, sorter []string, done chan struct{}) AccumulatorChannels{
	tmp := AccumulatorChannels{
		AggSrc: UpdateFeatureSpace2(subspace_channel_container.AggSrc, sorter, done),
		AggDst: UpdateFeatureSpace2(subspace_channel_container.AggDst, sorter, done),
		AggSrcDst: UpdateFeatureSpace2(subspace_channel_container.AggSrcDst, sorter, done),
	}
	return tmp
}

func UpdateFeatureSpace2(subspace_channels anomalies.SubspaceChannels, sorter []string, done chan struct{}) AccumulatorChannel{
	Xs := []preprocess.MicroSlot{}
	acc_c := make(chan preprocess.MicroSlot)

	go func() {
	LOOP:
		for {
			select {
			case X := <-acc_c:
				fmt.Print(".")

				if len(Xs) < preprocess.WINDOW_ARR_LEN-1 {
					Xs = append(Xs, X)
				} else {
					Xs = append(Xs, X)
					var x_old, x_new_update []tree.Point

					if len(Xs) == preprocess.WINDOW_ARR_LEN {
						//fmt.Println("before flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
						//fmt.Println("before flow processing")
						x_old = []tree.Point{}

						tmp := preprocess.MicroSlot{}
						for _, X := range Xs {
							tmp = append(tmp, X...)
						}
						x_new_update = preprocess.Normalize(tmp, sorter)
					} else if len(Xs) > preprocess.WINDOW_ARR_LEN {
						//fmt.Println("flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
						//fmt.Println("flow processing")
						x_old = Xs[0]
						Xs = Xs[1:]

						tmp := preprocess.MicroSlot{}
						for _, X := range Xs {
							tmp = append(tmp, X...)
						}
						x_new_update = preprocess.Normalize(tmp, sorter)
					}

					tmp1_x_old := map[[2]string][]tree.Point{}
					tmp2_x_new_update := map[[2]string][]tree.Point{}

					for _, p := range x_old {
						for subspace_keys, _ := range subspace_channels {
							subspace_key_0 := subspace_keys[0]
							subspace_key_1 := subspace_keys[1]
							subspace_val_0 := p.Vec_map[subspace_key_0]
							subspace_val_1 := p.Vec_map[subspace_key_1]

							point := tree.Point{
								Id:  p.Id,
								Key: p.Key,
								Vec: []float64{subspace_val_0, subspace_val_1},
								Vec_map: map[string]float64{
									subspace_key_0: subspace_val_0,
									subspace_key_1: subspace_val_1,
								},
							}

							tmp1_x_old[subspace_keys] = append(tmp1_x_old[subspace_keys], point)
						}
					}

					for _, p := range x_new_update {
						for subspace_keys, _ := range subspace_channels {
							subspace_key_0 := subspace_keys[0]
							subspace_key_1 := subspace_keys[1]
							subspace_val_0 := p.Vec_map[subspace_key_0]
							subspace_val_1 := p.Vec_map[subspace_key_1]

							point := tree.Point{
								Id:  p.Id,
								Key: p.Key,
								Vec: []float64{subspace_val_0, subspace_val_1},
								Vec_map: map[string]float64{
									subspace_key_0: subspace_val_0,
									subspace_key_1: subspace_val_1,
								},
							}

							tmp2_x_new_update[subspace_keys] = append(tmp1_x_old[subspace_keys], point)
						}
					}

					for subspace_keys, channel := range subspace_channels {
						anom := anomalies.ProcessPackage{
							X_old:        tmp1_x_old[subspace_keys],
							X_new_update: tmp2_x_new_update[subspace_keys],
						}
						channel <- anom
					}
				}
			case <-done:
				break LOOP
			default:
			}
		}
	}()
	return acc_c
}

func UpdateFeatureSpace(acc_c chan preprocess.X_micro_slot, data chan server.HttpData,sorter []string, subspaces map[[2]string]tree.Subspace, config Config){
	l := utils.Logger{}
	logger := l.New()
	count := 0
	Xs := []preprocess.X_micro_slot{}

	for{
		select{
		case X := <- acc_c:
			fmt.Print(".")

			if len(Xs) < preprocess.WINDOW_ARR_LEN - 1{
				Xs = append(Xs, X)
			} else {
				Xs = append(Xs, X)
				var x_old, x_new_update []tree.Point

				if len(Xs) == preprocess.WINDOW_ARR_LEN{
					//fmt.Println("before flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					//fmt.Println("before flow processing")
					var aggdst, aggsrc, aggsrcdst []tree.Point
					for _, X := range Xs{
						aggdst = append(aggdst, X.AggDst...)
						aggsrc = append(aggsrc, X.AggSrc...)
						aggsrcdst = append(aggsrcdst, X.AggSrcDst...)
					}
					x_old = []tree.Point{}
					x_new_update = preprocess.Normalize(aggdst, sorter)
				} else if len(Xs) > preprocess.WINDOW_ARR_LEN{
					//fmt.Println("flow processing data", preprocess.WINDOW_ARR_LEN, len(base_matrix))
					//fmt.Println("flow processing")
					Xs_old := Xs[0]
					var aggdst_old, aggsrc_old, aggsrcdst_old []tree.Point
					aggdst_old = append(aggdst_old, Xs_old.AggDst...)
					aggsrc_old = append(aggsrc_old, Xs_old.AggSrc...)
					aggsrcdst_old = append(aggsrcdst_old, X.AggSrcDst...)

					Xs = Xs[1:]
					var aggdst, aggsrc, aggsrcdst []tree.Point
					for _, X := range Xs{
						aggdst = append(aggdst, X.AggDst...)
						aggsrc = append(aggsrc, X.AggSrc...)
						aggsrcdst = append(aggsrcdst, X.AggSrcDst...)
					}

					x_old = aggdst_old
					x_new_update = preprocess.Normalize(aggdst, sorter)
				}

				m := map[[2]string]tree.Subspace{}
				if (config.Execution_type == SEQUENTIAL){
					m = SequentialClustering(subspaces, config, x_old, x_new_update)

				} else if (config.Execution_type == PARALLEL){
					num_clusterer := runtime.GOMAXPROCS(config.Num_cpu) //gets the current number of cores
					func (){
						defer utils.TimeTrack(time.Now(),  "Clustering", num_clusterer, logger)
						//ParallelClustering(num_clusterer, subspaces, config, x_old, x_new_update)
						fmt.Println("parallel clustering started", "x_old", len(x_old), "x_new_update", len(x_new_update))
						m = ParallelClustering(num_clusterer, subspaces, config, x_old, x_new_update)
						count++
						return
					}()
				}

				//http_data := ProcessDataForVisualization(subspaces)
				//data <- http_data
			}
		}
	}
}

func SequentialClustering(subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point)  map[[2]string]tree.Subspace{
	for _, subspace := range subspaces{
		subspace.ComputeSubspace(x_old, x_new_update)
		subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
	}
	return subspaces
}

type result struct{
	subspace_key [2]string
	subspace tree.Subspace
}

type processPackage struct{
	subspace_key [2]string
	subspace tree.Subspace
	config Config
	x_old []tree.Point
	x_new_update []tree.Point
}

func Clusterer(done <-chan struct{}, processPackages <-chan processPackage , c chan<- result){
	for processPackage := range processPackages{
		subspace := processPackage.subspace
		config := processPackage.config
		x_old := processPackage.x_old
		x_new_update := processPackage.x_new_update
		subspace.ComputeSubspace(x_old, x_new_update)
		subspace.Cluster(config.Min_dense_points, config.Min_cluster_points)
		//dissimilarity_map := ComputeDissmilarityVector(subspace)
		if len(subspace.GetOutliers()) > 0{
			fmt.Println("key:",subspace.Subspace_key, "outliers:", subspace.GetOutliers(), "clusters:", subspace.GetClusters())
		}
		//TODO: COMPUTE ANOMALIES HERE
		select {
		case c <- result{processPackage.subspace_key, subspace}:
		case <- done:
			return
		}
	}
}

func SubspaceIterator(done <- chan struct{}, subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point) (<-chan processPackage){
	processPackages := make(chan processPackage)
	go func(){
		//fmt.Println("Subspace len:", len(subspaces))
		defer close(processPackages)

		for _, subspace := range subspaces {
			processPackage := processPackage{
				subspace_key: subspace.Subspace_key,
				subspace: subspace,
				config: config,
				x_old: x_old,
				x_new_update: x_new_update,
			}
			select {
			case processPackages <- processPackage:
			case <-done:
				break
			}
		}
		return
	}()
	return processPackages
}

func ParallelClustering(num_clusterers int, subspaces map[[2]string]tree.Subspace, config Config, x_old []tree.Point, x_new_update []tree.Point)  map[[2]string]tree.Subspace{
	done := make(chan struct{})
	defer close(done)

	processPackages := SubspaceIterator(done, subspaces, config ,x_old, x_new_update)

	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(num_clusterers)
	for i := 0; i < num_clusterers; i++{
		go func(){
			Clusterer(done, processPackages, c)
			wg.Done()
		}()
	}
	go func(){
		wg.Wait()
		close(c)
	}()

	m := make(map[[2]string]tree.Subspace)
	for r := range c{
		//to extract needed data for server visualization
		m[r.subspace_key] = r.subspace
	}

	return m
}