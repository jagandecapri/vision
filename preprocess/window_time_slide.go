package preprocess

import (
	"time"
	"log"
	"github.com/jagandecapri/vision/cmd"
)

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
						log.Println("Initialize Time")
						time_counter = packet_time
						acc = NewAccumulator()
					}

					if !time_counter.IsZero() && packet_time.After(time_counter.Add(cmd.DeltaTClusterData)){
						log.Println("packet_time > time_counter")
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

func simulator(acc_c_receive AccumulatorChannel, acc_c_send AccumulatorChannel, delta_t time.Duration){
	ticker := time.NewTicker(delta_t)
	tmp_c := make(chan MicroSlot, 2)

	go func(){
		defer func(){
			log.Println("close buffered channel read by ticker")
			close(tmp_c)
		}()
		for{
			select{
			case pts, open := <-acc_c_receive:
				if open{
					log.Println("acc_c_receive value received by simulator")
					tmp_c <- pts
				} else{
					return
				}
			default:
			}
		}
	}()

	go func(){
		defer func(){
			log.Println("close ticker")
			log.Println("close acc_c_send sending to update_feature_space channel")
			close(acc_c_send)
			ticker.Stop()
		}()

		for{
			select{
			case <-ticker.C:
				select{
				case pts_from_buffer, open := <-tmp_c:
					if open {
						log.Println("value sent to update_feature_space from simulator")
						acc_c_send <- pts_from_buffer
					} else {
						return
					}
				default:
				}
			default:
			}
		}
	}()
}

func WindowTimeSlideSimulator(acc_c_receive AccumulatorChannels, acc_c_send AccumulatorChannels, delta_t time.Duration){
	go simulator(acc_c_receive.AggSrc, acc_c_send.AggSrc, delta_t)
	go simulator(acc_c_receive.AggDst, acc_c_send.AggDst, delta_t)
	go simulator(acc_c_receive.AggSrcDst, acc_c_send.AggSrcDst, delta_t)
}
