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

	num_channels := len(d.Channels)
	EvidenceAccummulationForOutliers("network_scan_syn", store, num_channels, done, wg_channels)
}

func NewNetworkScanSYN() *NetworkScanSYN{
	return &NetworkScanSYN{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"perSYN", "nbDstPort"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "nbDsts"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "avgPktSize"}: make(chan DissimilarityVectorContainer),
	}}
}
