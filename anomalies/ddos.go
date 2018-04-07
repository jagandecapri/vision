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
	store := NewDissimilarityMapContainer()
	done := make(chan struct{})

	go func(chan struct{}) {
		done_counter := 0
		defer func(){
			close(done)
			log.Println("close ddos out channel")
		}()

		for{
			select{
				case dis_vector, open := <-d.Channels[[2]string{"nbSrcs", "avgPktSize"}]:
					if open{
						store.Store(dis_vector.Id, dis_vector)
					} else {
						done_counter++
					}
				case dis_vector, open := <-d.Channels[[2]string{"perICMP", "perSYN"}]:
					if open{
						store.Store(dis_vector.Id, dis_vector)
					} else {
						done_counter++
					}
				case dis_vector, open := <-d.Channels[[2]string{"nbSrcPort", "perICMP"}]:
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

	go func(done chan struct{}, wg_channels *sync.WaitGroup){
		defer func(){
			log.Println("closing iterating channel in ddos")
			wg_channels.Done()
		}()

		for{
			out := store.IterateDissimilarityMapContainer()
			for dmp := range out{
				if len(dmp.Dis_vector) == len(d.Channels){
					log.Println("All subspaces processed in ddos disimilarity vector")
					//TODO: Sort and Calculate Knee here, http_data sending

					tmp := make(map[int][]DissimilarityVector)

					for _, dissimilarity_vector_container := range dmp.Dis_vector{
						for _, dissimilarity_vector := range dissimilarity_vector_container.DissimilarityVectors{
							tmp[dissimilarity_vector.Id] = append(tmp[dissimilarity_vector.Id], dissimilarity_vector)
						}
					}

					log.Println("Data for knee sort in NetworkScanSYN ", tmp)
					//knee_data := make([]float64, len(tmp))
					//
					//for _, distance := range tmp{
					//	knee_data[]
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
	}(done, wg_channels)
}

func NewDDOS() *DDOS{
	return &DDOS{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"nbSrcs", "avgPktSize"}: make(chan DissimilarityVectorContainer),
		[2]string{"perICMP", "perSYN"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbSrcPort", "perICMP"}: make(chan DissimilarityVectorContainer),
	}}
}