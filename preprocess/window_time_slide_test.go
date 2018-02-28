package preprocess

import (
	"testing"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestWindowTimeSlide2(t *testing.T) {
	ch := make(chan PacketData)
	acc_c := make(chan X_micro_slot)

	go WindowTimeSlide(ch, acc_c)

	go func(){
		handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")

		if(err != nil){
			log.Fatal(err)
		}

		count := 0

		for {
			if count == 10000{
				log.Println("counter reached")
				close(ch)
				break
			}
			data, ci, err := handleRead.ReadPacketData()
			if err != nil && err != io.EOF {
				close(ch)
				log.Fatal(err)
			} else if err == io.EOF {
				close(ch)
				break
			} else {
				packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
				ch <- PacketData{Data: packet, Metadata: ci}
			}
			count++
		}
	}()


	for {
		_, open := <-acc_c
		if open == false{
			break
		}
	}

}
