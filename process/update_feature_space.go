package process

import (
	"github.com/jagandecapri/vision/preprocess"
	"github.com/jagandecapri/vision/tree"
	"fmt"
	"github.com/jagandecapri/vision/anomalies"
	"log"
	"github.com/jagandecapri/vision/cmd"
)

func UpdateFeatureSpaceBuilder(subspace_channel_container anomalies.SubspaceChannelsContainer, sorter []string) preprocess.AccumulatorChannels{
	tmp := preprocess.AccumulatorChannels{
		AggSrc:    UpdateFeatureSpace(subspace_channel_container.AggSrc, sorter, "agg_src"),
		AggDst:    UpdateFeatureSpace(subspace_channel_container.AggDst, sorter, "agg_dst"),
		AggSrcDst: UpdateFeatureSpace(subspace_channel_container.AggSrcDst, sorter, "agg_srcdst"),
	}
	return tmp
}

func UpdateFeatureSpace(subspace_channels anomalies.SubspaceChannels, sorter []string, agg_key string) preprocess.AccumulatorChannel{
	Xs := []preprocess.MicroSlot{}
	acc_c := make(chan preprocess.MicroSlot)

	go func() {
		defer func(){
			log.Println("close subspace channels")
			for _, channel := range subspace_channels {
				close(channel)
			}
		}()

		for {
			select {
			case X, open := <-acc_c:
				if open{
					fmt.Print(".")

					if len(Xs) < cmd.WindowArrayLen-1 {
						Xs = append(Xs, X)
					} else {
						Xs = append(Xs, X)
						var x_old, x_new_update []tree.Point

						if len(Xs) == cmd.WindowArrayLen {
							//log.Println("before flow processing data", cmd.WindowArrayLen)
							//log.Println("before flow processing")
							x_old = []tree.Point{}

							tmp := preprocess.MicroSlot{}
							for _, X := range Xs {
								tmp = append(tmp, X...)
							}
							//log.Println(len(Xs), len(Xs[0]), len(X))
							x_new_update = preprocess.Normalize(tmp, sorter)

							//DEBUGGING
							//log.Println(agg_key)
							//for _, pt := range x_new_update{
							//	log.Println(pt.Id, " ", pt.Vec, pt.Vec_map)
							//}
						} else if len(Xs) > cmd.WindowArrayLen {
							//log.Println("flow processing data", cmd.WindowArrayLen)
							//log.Println("flow processing")
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

								tmp2_x_new_update[subspace_keys] = append(tmp2_x_new_update[subspace_keys], point)
							}
						}

						for subspace_keys, channel := range subspace_channels {
							//log.Println("agg_key: ", agg_key, " subspace_keys: ", subspace_keys, " x_old: ", tmp1_x_old[subspace_keys], " x_new_update: ", tmp2_x_new_update[subspace_keys])
							anom := anomalies.ProcessPackage{
								X_old:        tmp1_x_old[subspace_keys],
								X_new_update: tmp2_x_new_update[subspace_keys],
							}
							channel <- anom
						}
					}
				} else {
					return
				}
			default:
			}
		}
	}()
	return acc_c
}