package preprocess

import (
	"time"
	"log"
)


var delta_t = 300 * time.Millisecond
var window = 15 * time.Second
var WINDOW_ARR_LEN = int(window.Seconds()/delta_t.Seconds())
var Point_ctr = 0

func WindowTimeSlide(ch chan PacketData, acc_c chan X_micro_slot, done chan struct{}){

	go func(){
		acc := NewAccumulator()
		time_init := time.Now()
		time_counter := time.Time{}

		LOOP:
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
						acc_c <- X
						log.Println("Time to read data:", time.Since(time_init))
						time_init = time.Now()
						time_counter = time.Time{}
					}

					acc.AddPacket(pd.Data)
				case <-done:
					break LOOP
				default:
			}
		}
	}()
}
