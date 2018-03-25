package anomalies

import (
	"sync"
	"log"
)

type DDOS struct{
	Channels map[[2]string] chan DissimilarityVectorContainer
}

func (d *DDOS) GetChannel(subspace_key [2]string) chan DissimilarityVectorContainer {
	return d.Channels[subspace_key]
}

func (d *DDOS) WaitOnChannels(wg_channels *sync.WaitGroup){
	out := make(chan DissimilarityVectorContainer)

	go func() {
		done_counter := 0
		defer func(){
			log.Println("close ddos out channel")
			close(out)
		}()

		for{
			select{
				case dis_vector, open := <-d.Channels[[2]string{"nbSrcs", "avgPktSize"}]:
					if open{
						out <- dis_vector
					} else {
						done_counter++
					}
				case dis_vector, open := <-d.Channels[[2]string{"perICMP", "perSYN"}]:
					if open{
						out <- dis_vector
					} else {
						done_counter++
					}
				case dis_vector, open := <-d.Channels[[2]string{"nbSrcPort", "perICMP"}]:
					if open{
						out <- dis_vector
					} else {
						done_counter++
					}
				default:
			}

			if done_counter == len(d.Channels){
				return
			}
		}
	}()

	go func(){
		store := map[int][]DissimilarityVectorContainer{}
		defer func() {
			log.Println("Signal ddos waitgroup done")
			wg_channels.Done()
		}()

		for{
			select{
				case dis_vector, open := <-out:
					if open{
						store[dis_vector.Id] = append(store[dis_vector.Id], dis_vector)
						if len(store[dis_vector.Id]) == len(d.Channels){
							//TODO: Sort and Calculate Knee here, http_data sending
							//data_sort := []float64{}
							//for _, dissimilarity := range dissimilarity_map{
							//	data_sort = append(data_sort, dissimilarity)
							//}
							//
							//sort.Float64s(data_sort)
							//
							//kneedle := Kneedle{}
							//
							//if len(data_sort) > 0{
							//	knee := kneedle.Run(data_sort, 1, false)
							//	fmt.Println("data sort:", data_sort)
							//	fmt.Println("knee:",knee)
							//	if len(knee) > 0{
							//		for point_id, dissimilarity := range dissimilarity_map{
							//
							//		}
							//	}
							//}
							//http_data := ProcessDataForVisualization(subspaces)
							//data <- http_data
						}
					} else {
						return
					}
				default:
			}
		}
	}()
}

func NewDDOS() *DDOS{
	return &DDOS{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"nbSrcs", "avgPktSize"}: make(chan DissimilarityVectorContainer),
		[2]string{"perICMP", "perSYN"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbSrcPort", "perICMP"}: make(chan DissimilarityVectorContainer),
	}}
}