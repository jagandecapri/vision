package anomalies

import (
	"sync"
	"log"
)

type NetworkScanSYN struct{
	Channels map[[2]string] chan DissimilarityVectorContainer
}

func (d *NetworkScanSYN) GetChannel(subspace_key [2]string) chan DissimilarityVectorContainer {
	return d.Channels[subspace_key]
}

func (d *NetworkScanSYN) WaitOnChannels(wg_channels *sync.WaitGroup){
	store := NewDissimilarityMapContainer()
	done := make(chan struct{})

	go func(chan struct{}) {
		done_counter := 0
		defer func(){
			close(done)
			log.Println("close network syn out channel")
		}()

		for{
			select{
			case dis_vector, open := <-d.Channels[[2]string{"perSYN", "nbDstPort"}]:
				if open{
					store.Store(dis_vector.Id, dis_vector)
				} else {
					done_counter++
				}
			case dis_vector, open := <-d.Channels[[2]string{"nbDstPort", "nbDsts"}]:
				if open{
					store.Store(dis_vector.Id, dis_vector)
				} else {
					done_counter++
				}
			case dis_vector, open := <-d.Channels[[2]string{"nbDstPort", "avgPktSize"}]:
				if open{
					store.Store(dis_vector.Id, dis_vector)
				} else {
					done_counter++
				}
			default:
			}

			if done_counter == len(d.Channels){
				return
			}
		}
	}(done)

	go func(chan struct{}){
		for{
			out := store.IterateDissimilarityMapContainer()
			for dmp := range out{
				if len(dmp.Dis_vector) == len(d.Channels){
					log.Println("All subspaces processed in disimilarity vector")
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
					store.Delete(dmp.Key)
				}
			}

			select{
			case <-done:
				return
			default:
			}
		}
	}(done)
}

func NewNetworkScanSYN() *NetworkScanSYN{
	return &NetworkScanSYN{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"perSYN", "nbDstPort"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "nbDsts"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "avgPktSize"}: make(chan DissimilarityVectorContainer),
	}}
}
