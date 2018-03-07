package preprocess

import (
	"time"
	"log"
)


var delta_t = 300 * time.Millisecond
var window = 15 * time.Second
var WINDOW_ARR_LEN = int(window.Seconds()/delta_t.Seconds())
var Point_ctr = 0

type AccumulatorChannel chan MicroSlot

type AccumulatorChannels struct{
	AggSrc AccumulatorChannel
	AggDst AccumulatorChannel
	AggSrcDst AccumulatorChannel
}

func WindowTimeSlide(ch chan PacketData, acc_c AccumulatorChannels, done chan struct{}){

	go func(){
		acc := NewAccumulator()
		time_init := time.Now()
		time_counter := time.Time{}

		for{
			select{
				case pd := <- ch:
					packet_time := pd.Metadata.Timestamp

					if time_counter.IsZero(){
						//fmt.Println("Initialize Time")
						time_counter = packet_time
						acc = NewAccumulator()
					}

					if !time_counter.IsZero() && packet_time.After(time_counter.Add(delta_t)){
						//fmt.Println("packet_time > time_counter")
						X := acc.GetMicroSlot()
						acc_c.AggSrc <- X.AggSrc
						acc_c.AggDst <- X.AggDst
						acc_c.AggSrcDst <- X.AggSrcDst
						log.Println("Time to read data:", time.Since(time_init))
						time_init = time.Now()
						time_counter = time.Time{}
					}

					acc.AddPacket(pd.Data)
				case <-done:
					return
				default:
			}
		}
	}()
}
