package anomalies

type NetworkScanSYN struct{
	Channels map[[2]string] chan DissimilarityVectorContainer
}

func (d *NetworkScanSYN) GetChannel(subspace_key [2]string) chan DissimilarityVectorContainer {
	return d.Channels[subspace_key]
}

func (d *NetworkScanSYN) WaitOnChannels(done chan struct{}){
	out := make(chan DissimilarityVectorContainer)
	go func() {
	LOOP:
		for{
			select{
			case dis_vector := <-d.Channels[[2]string{"perSYN", "nbDstPort"}]:
				out <- dis_vector
			case dis_vector := <-d.Channels[[2]string{"nbDstPort", "nbDsts"}]:
				out <- dis_vector
			case dis_vector := <-d.Channels[[2]string{"nbDstPort", "avgPktSize"}]:
				out <- dis_vector
			case <-done:
				break LOOP
			default:
			}
		}
	}()

	go func(){
		store := map[int][]DissimilarityVectorContainer{}
	LOOP:
		for{
			select{
			case dis_vector := <-out:
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
			case <-done:
				break LOOP
			default:
			}
		}
	}()
}

func NewNetworkScanSYN() *NetworkScanSYN{
	return &NetworkScanSYN{Channels: map[[2]string] chan DissimilarityVectorContainer{
		[2]string{"perSYN", "nbDstPort"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "nbDsts"}: make(chan DissimilarityVectorContainer),
		[2]string{"nbDstPort", "avgPktSize"}: make(chan DissimilarityVectorContainer),
	}}
}
