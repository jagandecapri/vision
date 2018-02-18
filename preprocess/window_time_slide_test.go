package preprocess

import (
	"testing"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math"
)

func TestWindowTimeSlide2(t *testing.T) {
	ch := make(chan PacketData)
	acc := make(chan Accumulator)
	quit := make(chan int)

	go WindowTimeSlide2(ch, acc, quit)

	go func(){
		handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")

		if(err != nil){
			log.Fatal(err)
		}

		count := 0

		for {
			if count == math.MaxInt64{
				quit <- 0
				break
			}
			data, ci, err := handleRead.ReadPacketData()
			if err != nil && err != io.EOF {
				quit <- 0
				log.Fatal(err)
			} else if err == io.EOF {
				quit <- 0
				break
			} else {
				packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
				ch <- PacketData{packet, ci}
			}
			count++
		}
	}()


	for {
		tmp_acc, open := <-acc
		if open == false{
			break
		}
	}

}
