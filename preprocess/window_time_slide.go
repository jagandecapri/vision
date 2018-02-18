package preprocess

import (
	"fmt"
	"time"
	"github.com/google/gopacket"
)


var delta_t = 300 * time.Millisecond
var window = 15 * time.Second
var WINDOW_ARR_LEN = int(window.Seconds()/delta_t.Seconds())

func WindowTimeSlide(ch chan PacketData, acc chan PacketAcc, quit chan int){
	time_counter := time.Time{}
	tmp := []PacketFeature{}
	initializePacketAcc := func(){
		tmp = []PacketFeature{}
	}
	for{
		select{
		case pd := <- ch:
			packet_time := pd.Metadata.Timestamp
			if time_counter.IsZero() {
				//fmt.Println("Initialize Time")
				time_counter = packet_time
				pf := PacketFeature{}
				pf.ExtractFeature(pd.Data)
				initializePacketAcc()
				tmp = append(tmp, pf)
			} else if packet_time.Before(time_counter.Add(delta_t)) || packet_time.Equal(time_counter.Add(delta_t)) {
				//fmt.Println("packet_time <= time_counter")
				pf := PacketFeature{}
				pf.ExtractFeature(pd.Data)
				tmp = append(tmp, pf)
			} else if packet_time.After(time_counter.Add(delta_t)) {
				//fmt.Println("packet_time > time_counter")
				acc <- tmp
				initializePacketAcc()
				time_counter = time_counter.Add(delta_t)
			}
		case <-quit:
			fmt.Println("No Data")
			close(ch)
			close(acc)
		}
	}
}

type Accumulator struct{
	AggSrc map[gopacket.Endpoint][]PacketFeature
	AggDst map[gopacket.Endpoint][]PacketFeature
	AggSrcDst map[gopacket.Flow][]PacketFeature
}

func NewAccumulator() Accumulator{
	return Accumulator{AggSrc: make(map[gopacket.Endpoint][]PacketFeature),
		AggDst: make(map[gopacket.Endpoint][]PacketFeature),
		AggSrcDst: make(map[gopacket.Flow][]PacketFeature),
	}
}

func WindowTimeSlide2(ch chan PacketData, acc chan Accumulator, quit chan int){
	time_counter := time.Time{}
	var tmp_acc Accumulator

	for{
		select{
		case pd, open := <- ch:
			if open == true{
				packet_time := pd.Metadata.Timestamp

				if time_counter.IsZero() {
					fmt.Println("Initialize Time")
					time_counter = packet_time
					pf := PacketFeature{}
					pf.ExtractFeature(pd.Data)

					tmp_acc = NewAccumulator()
					netFlow := pd.Data.NetworkLayer().NetworkFlow()
					tmp_acc.AggSrc[netFlow.Src()] = append(tmp_acc.AggSrc[netFlow.Src()], pf)
					tmp_acc.AggDst[netFlow.Dst()] = append(tmp_acc.AggDst[netFlow.Dst()], pf)
					tmp_acc.AggSrcDst[netFlow] = append(tmp_acc.AggSrc[netFlow.Dst()], pf)
				} else if packet_time.Before(time_counter.Add(delta_t)) || packet_time.Equal(time_counter.Add(delta_t)) {
					//fmt.Println("packet_time <= time_counter")
					pf := PacketFeature{}
					pf.ExtractFeature(pd.Data)
					netFlow := pd.Data.NetworkLayer().NetworkFlow()
					tmp_acc.AggSrc[netFlow.Src()] = append(tmp_acc.AggSrc[netFlow.Src()], pf)
					tmp_acc.AggDst[netFlow.Dst()] = append(tmp_acc.AggDst[netFlow.Dst()], pf)
					tmp_acc.AggSrcDst[netFlow] = append(tmp_acc.AggSrc[netFlow.Dst()], pf)
				} else if packet_time.After(time_counter.Add(delta_t)) {
					fmt.Println("packet_time > time_counter")
					acc <- tmp_acc
					tmp_acc = NewAccumulator()
					time_counter = time_counter.Add(delta_t)
				}
			}
		case <-quit:
			fmt.Println("No Data")
			close(ch)
			close(acc)
		}
	}
}
