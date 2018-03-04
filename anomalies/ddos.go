package anomalies

import (
	"github.com/jagandecapri/vision/process"
)

type DDOS struct{
	Channels map[[2]string] chan process.DissimilarityVector
}

func (d *DDOS) GetChannel(subspace_key [2]string) chan process.DissimilarityVector {
	return d.Channels[subspace_key]
}

func (d *DDOS) WaitOnChannels() chan bool {
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

	done := make(chan bool)
	go func() {
		LOOP:
		for{
			select{
				case <-d.Channels[[2]string{"nbSrcs", "avgPktSize"}]:
				case <-d.Channels[[2]string{"perICMP", "perSYN"}]:
				case <-d.Channels[[2]string{"nbSrcPort", "perICMP"}]:
				case <- done:
					break LOOP
			}
		}
	}()
	return done
}

func NewDDOS() *DDOS{
	return &DDOS{Channels: map[[2]string] chan process.DissimilarityVector{
		[2]string{"nbSrcs", "avgPktSize"}: make(chan process.DissimilarityVector),
		[2]string{"perICMP", "perSYN"}: make(chan process.DissimilarityVector),
		[2]string{"nbSrcPort", "perICMP"}: make(chan process.DissimilarityVector),
	}}
}