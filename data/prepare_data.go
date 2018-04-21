package data

import (
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jagandecapri/vision/preprocess"
	"time"
	"fmt"
)

var Point_ctr = 0

func WindowTimeSlide(ch chan preprocess.PacketData, acc_c preprocess.AccumulatorChannels, delta_t time.Duration, done chan struct{}){

	go func(){
		acc := preprocess.NewAccumulator()
		//time_init := time.Now()
		time_counter := time.Time{}

		for{
			select{
			case pd := <- ch:
				fmt.Printf("*")
				packet_time := pd.Metadata.Timestamp

				if time_counter.IsZero(){
					//log.Println("Initialize Time")
					time_counter = packet_time
					acc = preprocess.NewAccumulator()
				}

				if !time_counter.IsZero() && packet_time.After(time_counter.Add(delta_t)){
					//log.Println("packet_time > time_counter")
					X := acc.GetMicroSlot()
					acc_c.AggSrc <- X.AggSrc
					acc_c.AggDst <- X.AggDst
					acc_c.AggSrcDst <- X.AggSrcDst
					//log.Println("Time to read data:", time.Since(time_init))
					//time_init = time.Now()
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

func Run(pcap_file_path string, db_name string, delta_t time.Duration){
	ch := make(chan preprocess.PacketData)
	done := make(chan struct{})
	acc_c := preprocess.AccumulatorChannels{
		AggSrc:    make(chan preprocess.MicroSlot),
		AggDst:    make(chan preprocess.MicroSlot),
		AggSrcDst: make(chan preprocess.MicroSlot),
	}


	WindowTimeSlide(ch, acc_c, delta_t, done)
	NewSQL(db_name, acc_c, done, delta_t)
	handleRead, err := pcap.OpenOffline(pcap_file_path)

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			close(done)
			log.Fatal(err)
		} else if err == io.EOF {
			close(done)
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- preprocess.PacketData{Data: packet, Metadata: ci}
		}
	}
}
