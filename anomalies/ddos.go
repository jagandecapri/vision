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

	num_channels := len(d.Channels)
	EvidenceAccummulationForOutliers("ddos", store, num_channels, done, wg_channels)
}

func NewDDOS() *DDOS{
	return &DDOS{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"nbSrcs", "avgPktSize"}: make(chan DissimilarityVectorContainer),
		[2]string{"perICMP", "perSYN"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbSrcPort", "perICMP"}: make(chan DissimilarityVectorContainer),
	}}
}