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
	go func() {
		done_counter := 0
		num_channels := len(d.Channels)
		log.Println("num_channels", num_channels)
		store := NewDissimilarityMapContainer2("network_scan_syn", num_channels)
		done := make(chan struct{})

		defer func(){
			close(done)
			wg_channels.Done()
			log.Println("close network_scan_syn out channel")
		}()

		for id, channel := range d.Channels{
			go func(id [2]string, store *DissimilarityMapContainer, channel chan DissimilarityVectorContainer){
				for{
					select{
					case dis_vector, open := <-channel:
						if open{
							store.Store(dis_vector.Id, dis_vector)
						} else {
							done<- struct{}{}
							return
						}
					}
				}
			}(id, store, channel)
		}

		for{
			select{
			case <-done:
				done_counter++
				if done_counter == len(d.Channels){
					return
				}
			}
		}
	}()
}

func NewNetworkScanSYN() *NetworkScanSYN {
	return &NetworkScanSYN{Channels: map[[2]string]chan DissimilarityVectorContainer{
		[2]string{"perSYN", "nbDstPort"}:     make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "nbDsts"}:     make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "avgPktSize"}: make(chan DissimilarityVectorContainer),
	}}
}