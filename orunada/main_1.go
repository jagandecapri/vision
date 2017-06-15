package main

import (
	"github.com/jagandecapri/vision/orunada/utils"
	"github.com/jagandecapri/vision/orunada/grid"
	"github.com/jagandecapri/vision/orunada/preprocess"
	"github.com/cockroachdb/apd"
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func getSorter() []string{
	sorter := []string{}
	sorter = append(sorter, "nbPacket", "nbSrcPort", "nbDstPort", "nbSrcs", "nbDsts", "perSyn", "perAck", "perRST", "perFIN", "perCWR", "perURG", "avgPktSize", "meanTTL")
	return sorter
}

func GetSubspaceKey(sorter []string, feature_cnt int) [][]string {
	all := [][]string{}
	utils.Comb(len(sorter), feature_cnt, func (c []int){
		tmp := []string{}
		for _, v := range c {
			tmp = append(tmp, sorter[v])
		}
		all = append(all, tmp)
	})
	return all
}

func main(){
	sorter:= getSorter()
	subspace_keys := GetSubspaceKey(sorter, 2)
	ctx := apd.BaseContext.WithPrecision(6)
	grid := new(grid.Grid)
	for _, subspace_key := range subspace_keys{
		g := grid.Build2DGrid(subspace_key, ctx)
		fmt.Println(g)
	}

	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan preprocess.PacketData)
	acc := make(chan preprocess.PacketAcc)
	quit := make(chan int)
	go preprocess.WindowTimeSlide(ch, acc, quit)
	//go UpdateFS(acc)

	if(err != nil){
		log.Fatal(err)
	}

	for {
		data, ci, err := handleRead.ReadPacketData()
		if err != nil && err != io.EOF {
			quit <- 0
			log.Fatal(err)
		} else if err == io.EOF {
			quit <- 0
			break
		} else {
			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			ch <- preprocess.PacketData{packet, ci}
		}
	}
}
