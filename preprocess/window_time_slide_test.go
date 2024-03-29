package preprocess

import (
	"testing"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestWindowTimeSlide(t *testing.T) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "C:\\Users\\Jack\\go\\src\\github.com\\jagandecapri\\vision\\logs\\lumber_log.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true, // disabled by default
	})

	ch := make(chan PacketData)
	acc_c := AccumulatorChannels{
		AggSrc: make(AccumulatorChannel),
		AggDst: make(AccumulatorChannel),
		AggSrcDst: make(AccumulatorChannel),
	}

	done := make(chan struct{})
	go WindowTimeSlide(ch, acc_c, done)

	go func(){
		handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")

		if(err != nil){
			log.Fatal(err)
		}

		count := 0

		for {
			if count == 100000{
				log.Println("counter reached")
				close(done)
				break
			}
			data, ci, err := handleRead.ReadPacketData()
			if err != nil && err != io.EOF {
				close(done)
				log.Fatal(err)
			} else if err == io.EOF {
				close(done)
				break
			} else {
				packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
				ch <- PacketData{Data: packet, Metadata: ci}
			}
			count++
		}
	}()


	for {
		select{
			case <-acc_c.AggSrc:
			case <-acc_c.AggDst:
			case <-acc_c.AggSrcDst:
			case <-done:
				return
			default:
		}
	}

}
