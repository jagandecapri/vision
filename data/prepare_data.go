package data

import (
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jagandecapri/vision/preprocess"
	"os"
	"time"
	"fmt"
)

var delta_t = 300 * time.Millisecond
var window = 15 * time.Second
var WINDOW_ARR_LEN = int(window.Seconds()/delta_t.Seconds())
var Point_ctr = 0

func WindowTimeSlide(ch chan preprocess.PacketData, acc_c preprocess.AccumulatorChannels, done chan struct{}){

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

func Run(){
	ch := make(chan preprocess.PacketData)
	done := make(chan struct{})
	acc_c := preprocess.AccumulatorChannels{
		AggSrc:    make(chan preprocess.MicroSlot),
		AggDst:    make(chan preprocess.MicroSlot),
		AggSrcDst: make(chan preprocess.MicroSlot),
	}

	pcap_file_path := os.Getenv("PCAP_FILE")
	if pcap_file_path == ""{
		pcap_file_path = "C:\\Users\\Jack\\Downloads\\201705021400.pcap"
	}

	WindowTimeSlide(ch, acc_c, done)
	NewSQL(acc_c, done, delta_t)
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
