package main

import (
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF {
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			// Iterate over all layers, printing out each layer type
			for _, layer := range packet.Layers() {
				fmt.Println(packet.String(), layer.LayerType(), ci)
			}
		}
	}

}