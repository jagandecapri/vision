package preprocess

import (
	"time"
	"log"
)


var delta_t = 300 * time.Millisecond
var window = 600 * time.Millisecond
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
						log.Println("Initialize Time")
						time_counter = packet_time
						acc = NewAccumulator()
					}

					if !time_counter.IsZero() && packet_time.After(time_counter.Add(delta_t)){
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

func WindowTimeSlideSimulator(acc_c_receive AccumulatorChannels, acc_c_send AccumulatorChannels, delta_t time.Duration, done chan struct{}){

	go func(acc_c_receive AccumulatorChannels, acc_c_send AccumulatorChannels, delta_t time.Duration){
		ticker := time.NewTicker(delta_t)
		tmp_c := make(chan MicroSlot, 2)
		for{
			select{
			case pts := <-acc_c_receive.AggSrc:
				tmp_c <- pts
				//log.Println("Received aggsrc data")
			case <-ticker.C:
				select{
					case pts_from_buffer := <-tmp_c:
						acc_c_send.AggSrc <- pts_from_buffer
					default:
				}
			case <-done:
				ticker.Stop()
				return
			default:
			}
		}
	}(acc_c_receive, acc_c_send, delta_t)

	go func(acc_c_receive AccumulatorChannels, acc_c_send AccumulatorChannels, delta_t time.Duration){
		ticker := time.NewTicker(delta_t)
		tmp_c := make(chan MicroSlot, 2)
		for{
			select{
			case pts := <-acc_c_receive.AggDst:
				tmp_c <- pts
				//log.Println("Received aggdst data")
			case <-ticker.C:
				select{
					case pts_from_buffer := <-tmp_c:
						acc_c_send.AggDst <- pts_from_buffer
					default:
				}
			case <-done:
				ticker.Stop()
				return
			default:
			}
		}
	}(acc_c_receive, acc_c_send, delta_t)

	go func(acc_c_receive AccumulatorChannels, acc_c_send AccumulatorChannels, delta_t time.Duration){
		ticker := time.NewTicker(delta_t)
		tmp_c := make(chan MicroSlot, 2)
		for{
			select{
			case pts := <-acc_c_receive.AggSrcDst:
				tmp_c <- pts
				//log.Println("Received aggsrcdst data")
			case <-ticker.C:
				select{
					case pts_from_buffer := <-tmp_c:
						acc_c_send.AggSrcDst <- pts_from_buffer
					default:
				}
			case <-done:
				ticker.Stop()
				return
			default:
			}
		}
	}(acc_c_receive, acc_c_send, delta_t)
}
