package process

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/tree"
	"fmt"
	"github.com/jagandecapri/vision/anomalies"
)

func UpdateFeatureSpaceBuilder(subspace_channel_container anomalies.SubspaceChannelsContainer, sorter []string, done chan struct{}) preprocess.AccumulatorChannels{
	tmp := preprocess.AccumulatorChannels{
		AggSrc:    UpdateFeatureSpace(subspace_channel_container.AggSrc, sorter, done),
		AggDst:    UpdateFeatureSpace(subspace_channel_container.AggDst, sorter, done),
		AggSrcDst: UpdateFeatureSpace(subspace_channel_container.AggSrcDst, sorter, done),
	}
	return tmp
}

func UpdateFeatureSpace(subspace_channels anomalies.SubspaceChannels, sorter []string, done chan struct{}) preprocess.AccumulatorChannel{
	Xs := []preprocess.MicroSlot{}
	acc_c := make(chan preprocess.MicroSlot)

	go func() {
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
				return
			default:
			}
		}
	}()
	return acc_c
}