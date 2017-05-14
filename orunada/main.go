package main

import (
	"github.com/google/gopacket/pcap"
	"log"
	"io"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
	"net"
)

type PacketFeature struct{
	SrcIP, DstIP net.IP
	SrcPort, DstPort layers.TCPPort
	SYN, ACK, RST, FIN, CWR, URG bool
	length uint16
	TTL uint8
}
func extractFeature(packet gopacket.Packet) PacketFeature{
	packet_feature := PacketFeature{}
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		// ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	}

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		packet_feature.SrcIP = ipv4.SrcIP
		packet_feature.DstIP = ipv4.DstIP
		packet_feature.TTL = ipv4.TTL
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		packet_feature.SrcIP = ipv6.SrcIP
		packet_feature.DstIP = ipv6.DstIP
		packet_feature.TTL = ipv6.TTL
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		packet_feature.SrcPort = tcp.SrcPort
		packet_feature.DstPort = tcp.DstPort
		packet_feature.SYN = tcp.SYN
		packet_feature.ACK = tcp.ACK
		packet_feature.RST = tcp.RST
		packet_feature.FIN = tcp.FIN
		packet_feature.CWR = tcp.CWR
		packet_feature.URG = tcp.URG
	}

	// Let's see if the packet is icmp4
	icmp4Layer := packet.Layer(layers.LayerTypeICMPv4)
	if icmp4Layer != nil {
		//fmt.Println("ICMP4 Layer Detected")
		//icmp4, _ := icmp4Layer.(*layers.ICMPv4)
		//fmt.Println(icmp4.TypeCode.String())
	}

	// Let's see if the packet is icmp6
	icmp6Layer := packet.Layer(layers.LayerTypeICMPv6)
	if icmp6Layer != nil {
		//fmt.Println("ICMP6 Layer Detected")
		//icmp6, _ := icmp6Layer.(*layers.ICMPv6)
		//fmt.Println(icmp6.TypeCode.String())
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		//fmt.Println("Error decoding some part of the packet:", err)
	}
	return packet_feature
}

var delta_t time.Duration = 300 * time.Millisecond
var window time.Duration = 15 * time.Second
var window_arr_len = int(window.Seconds()/delta_t.Seconds())
var time_counter time.Time

type PacketData struct{
	data gopacket.Packet
	metadata gopacket.CaptureInfo
}

type PacketAcc []PacketFeature

func UniqFloat64(input []float64) []float64 {
	u := make([]float64, 0, len(input))
	m := make(map[float64]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func UniqString(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func ExtractDeltaPacketFeature(feature_packets PacketAcc) []float64{
	nbSYN, nbACK, nbRST, nbFIN, nbCWR, nbURG, totalPktSize, totalTTL  := 0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0
	srcPorts := []float64{}
	dstPorts := []float64{}
	srcIPs := []string{}
	dstIPs := []string{}

	for _, fp := range feature_packets{
		if fp.SrcPort != 0 {
			srcPorts = append(srcPorts, float64(fp.SrcPort))
		}
		if fp.DstPort != 0 {
			dstPorts = append(dstPorts, float64(fp.DstPort))
		}
		if fp.SrcIP != nil {
			srcIPs = append(srcIPs, string(fp.SrcIP))
		}
		if fp.DstIP != nil {
			dstIPs = append(dstIPs, string(fp.DstIP))
		}
		if fp.SYN == true{
			nbSYN++
		}
		if fp.ACK == true{
			nbACK++
		}
		if fp.RST == true{
			nbRST++
		}
		if fp.ACK == true{
			nbFIN++
		}
		if fp.CWR == true{
			nbCWR++
		}
		if fp.URG == true{
			nbURG++
		}
		totalPktSize += float64(fp.length)
		totalTTL += float64(fp.TTL)
	}
	x := []float64{}
	nbPacket := float64(len(feature_packets))
	nbSrcPort := float64(len(UniqFloat64(srcPorts)))
	nbDstPort := float64(len(UniqFloat64(dstPorts)))
	nbSrcs := float64(len(UniqString(srcIPs)))
	nbDsts := float64(len(UniqString(dstIPs)))
	perSyn := nbSYN/nbPacket*100
	perAck := nbACK/nbPacket*100
	perRST := nbRST/nbPacket*100
	perFIN := nbFIN/nbPacket*100
	perCWR := nbCWR/nbPacket*100
	perURG := nbURG/nbPacket*100
	avgPktSize := totalPktSize/nbPacket
	meanTTL := totalTTL/nbPacket
	x = append(x, nbPacket, nbSrcPort, nbDstPort, nbSrcs, nbDsts, perSyn, perAck, perRST, perFIN, perCWR, perURG, avgPktSize, meanTTL)
	return x
}

func WindowTimeSlide(ch chan PacketData, acc chan PacketAcc, quit chan int){
	var tmp PacketAcc
	initializePacketAcc := func(){
		tmp = []PacketFeature{}
	}
	for{
		select{
		case pd := <- ch:
			packet_time := pd.metadata.Timestamp
			if time_counter.IsZero() {
				//fmt.Println("Initialize Time")
				time_counter = packet_time
				pf := extractFeature(pd.data)
				initializePacketAcc()
				tmp = append(tmp, pf)
			} else if packet_time.Before(time_counter.Add(delta_t)) || packet_time.Equal(time_counter.Add(delta_t)) {
				//fmt.Println("packet_time <= time_counter")
				pf := extractFeature(pd.data)
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

func UpdateFS(acc chan PacketAcc){
	arr := [][]float64{}
	for{
		select{
		case packet_acc := <- acc:
			fmt.Print(".")
			x := ExtractDeltaPacketFeature(packet_acc)
			if len(arr) != window_arr_len{
				arr = append(arr, x)
			} else {
				x_new := x
				/*x_old*/_, x_update := arr[0], arr[1:]
				arr = append(x_update, x_new)
				fmt.Println(arr)
			}

		}
	}
}

func main() {
	handleRead, err := pcap.OpenOffline("C:\\Users\\Jack\\Downloads\\201705021400.pcap")
	ch := make(chan PacketData)
	acc := make(chan PacketAcc)
	quit := make(chan int)
	go WindowTimeSlide(ch, acc, quit)
	go UpdateFS(acc)

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
			ch <- PacketData{packet, ci}
		}
	}

}