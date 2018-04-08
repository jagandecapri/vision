package anomalies

import (
	"sync"
	"log"
	"sort"
	"github.com/jagandecapri/vision/utils"
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

	go func(done chan struct{}, wg_channels *sync.WaitGroup){
		defer func(){
			log.Println("closing iterating channel in network scan syn")
			wg_channels.Done()
		}()

		counter := 0

		for{
			out := store.IterateDissimilarityMapContainer()
			for dmp := range out{
				if len(dmp.Dis_vector) == len(d.Channels){
					log.Println("All subspaces processed in disimilarity vector")

					tmp := make(map[int]DissimilarityVector)

					for _, dissimilarity_vector_container := range dmp.Dis_vector{
						for _, dissimilarity_vector := range dissimilarity_vector_container.DissimilarityVectors{
							if val, ok := tmp[dissimilarity_vector.Id]; ok{
								val.Distance += dissimilarity_vector.Distance
								tmp[dissimilarity_vector.Id] = val
							} else {
								tmp[dissimilarity_vector.Id] = dissimilarity_vector
							}
						}
					}

					type kv struct {
						Key   int
						Value float64
					}

					var ss []kv

					for k, v := range tmp {
						ss = append(ss, kv{Key: k, Value: v.Distance})
					}

					sort.Slice(ss, func(i, j int) bool {
						return ss[i].Value < ss[j].Value
					})

					var knee_data []float64

					for _, kv := range ss {
						knee_data = append(knee_data, kv.Value)
						val := tmp[kv.Key]
						srcIP := utils.UniqString(val.PointKey.SrcIP)
						dstIP := utils.UniqString(val.PointKey.DstIP)
						srcPort := utils.UniqString(val.PointKey.SrcPort)
						dstPort := utils.UniqString(val.PointKey.DstPort)

						log.Printf("Batch: %v Keys: SrcIP: %+v DstIP: %+v SrcPort: %+v DstPort: %+v Distance: %+v", counter, srcIP, dstIP, srcPort, dstPort, kv.Value)
					}

					kneedle := Kneedle{}
					if len(knee_data) > 0{
						knee := kneedle.Run(knee_data, 1, false)
						log.Println("network_scan_sync data sort: ", knee_data)
						log.Println("network_scan_sync knee: ",knee)
						//if len(knee) > 0{
						//	for point_id, dissimilarity := range dissimilarity_map{
						//
						//	}
						//}
					}
					store.Delete(dmp.Key)
					counter++
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

func NewNetworkScanSYN() *NetworkScanSYN{
	return &NetworkScanSYN{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"perSYN", "nbDstPort"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "nbDsts"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "avgPktSize"}: make(chan DissimilarityVectorContainer),
	}}
}
