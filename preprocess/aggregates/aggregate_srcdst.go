package aggregates

import (
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
	"net"
	"github.com/jagandecapri/vision/tree"
)

type AggSrcDst struct{
	FlowKeys []tree.PointKey
	nbPacket float64
	srcPort map[layers.TCPPort]int
	dstPort map[layers.TCPPort]int
	nbSYN float64
	nbACK float64
	nbRST float64
	nbFIN float64
	nbCWR float64
	nbURG float64
	totalPacketSize float64
	TTL float64
	nbICMP float64
	nbICMPRed float64
	nbICMPTime float64
	nbICMPUnr float64
	nbICMPOther float64
}

func (a *AggSrcDst) AddPacket(p gopacket.Packet) gopacket.ErrorLayer{

	a.nbPacket++

	// Let's see if the packet is an ethernet packet
	ethernetLayer := p.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		// ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	}

	var srcIP, dstIP net.IP
	var srcPort, dstPort layers.TCPPort

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		srcIP = ipv4.SrcIP
		dstIP = ipv4.DstIP
		a.TTL += float64(ipv4.TTL)
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		srcIP = ipv6.SrcIP
		dstIP = ipv6.DstIP
		a.TTL += float64(ipv6.TTL)
	}

	// Let's see if the packet is TCP
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = tcp.SrcPort
		dstPort = tcp.DstPort
		if tcp.SYN{
			a.nbSYN++
		}
		if tcp.ACK{
			a.nbACK++
		}
		if tcp.RST{
			a.nbRST++
		}
		if tcp.FIN{
			a.nbFIN++
		}
		if tcp.CWR{
			a.nbCWR++
		}
		if tcp.URG{
			a.nbURG++
		}
	}

	flowKey := tree.PointKey{SrcIP: srcIP, DstIP: dstIP, SrcPort: srcPort, DstPort: dstPort}
	a.FlowKeys = append(a.FlowKeys, flowKey)

	// Let's see if the packet is icmp4
	icmp4Layer :=p.Layer(layers.LayerTypeICMPv4)
	if icmp4Layer != nil {
		a.nbICMP++
		//fmt.Println("ICMP4 Layer Detected")
		icmp4, _ := icmp4Layer.(*layers.ICMPv4)
		//fmt.Println(icmp4.TypeCode.String())
		switch icmp4.TypeCode{
		case layers.ICMPv4TypeRedirect:
			a.nbICMPRed++
		case layers.ICMPv4TypeTimeExceeded:
			a.nbICMPTime++
		case layers.ICMPv4TypeDestinationUnreachable:
			a.nbICMPUnr++
		default:
			a.nbICMPOther++
		}
	}

	// Let's see if the packet is icmp6
	icmp6Layer := p.Layer(layers.LayerTypeICMPv6)
	if icmp6Layer != nil {
		a.nbICMP++
		//fmt.Println("ICMP6 Layer Detected")
		icmp6, _ := icmp6Layer.(*layers.ICMPv6)
		//fmt.Println(icmp6.TypeCode.String())
		switch icmp6.TypeCode{
		case layers.ICMPv6TypeRedirect:
			a.nbICMPRed++
		case layers.ICMPv6TypeTimeExceeded:
			a.nbICMPTime++
		case layers.ICMPv6TypeDestinationUnreachable:
			a.nbICMPUnr++
		default:
			a.nbICMPOther++
		}
	}

	// Check for errors
	var err gopacket.ErrorLayer
	if err = p.ErrorLayer(); err != nil {
		return err
	}
	return err
}

func (a *AggSrcDst) GetKey() []tree.PointKey {
	return a.FlowKeys
}

func (a *AggSrcDst) NbPacket() float64 {
	return a.nbPacket
}

func (a *AggSrcDst) NbSrcPort() float64 {
	return float64(len(a.srcPort))
}

func (a *AggSrcDst) NbDstPort() float64 {
	return float64(len(a.dstPort))
}

func (a *AggSrcDst) NbSrcs() float64 {
	return 0.0
}

func (a *AggSrcDst) NbDsts() float64 {
	return 0.0
}

func (a *AggSrcDst) PerSYN() float64 {
	return float64((a.nbSYN/a.nbPacket)*100)
}

func (a *AggSrcDst) PerACK() float64 {
	return float64((a.nbACK/a.nbPacket)*100)
}

func (a *AggSrcDst) PerRST() float64 {
	return float64((a.nbRST/a.nbPacket)*100)
}

func (a *AggSrcDst) PerFIN() float64 {
	return float64((a.nbFIN/a.nbPacket)*100)
}

func (a *AggSrcDst) PerCWR() float64 {
	return float64((a.nbCWR/a.nbPacket)*100)
}

func (a *AggSrcDst) PerURG() float64 {
	return float64((a.nbURG/a.nbPacket)*100)
}

func (a *AggSrcDst) AvgPktSize() float64 {
	return float64(a.totalPacketSize/a.nbPacket)
}

func (a *AggSrcDst) MeanTTL() float64 {
	return float64(a.TTL/a.nbPacket)
}

func (a *AggSrcDst) SimIPDst() float64 {
	return 0.0
}

func (a *AggSrcDst) SimIPSrc() float64 {
	return 0.0
}

func (a *AggSrcDst) PerICMP() float64{
	return float64((a.nbICMP/a.nbPacket)*100)
}

func (a *AggSrcDst) PerICMPRed() float64 {
	return float64((a.nbICMPRed/a.nbPacket)*100)
}

func (a *AggSrcDst) PerICMPTime() float64 {
	return float64((a.nbICMPTime/a.nbPacket)*100)
}

func (a *AggSrcDst) PerICMPUnr() float64 {
	return float64((a.nbICMPUnr/a.nbPacket)*100)
}

func (a *AggSrcDst) PerICMPOther() float64 {
	return float64((a.nbICMPOther/a.nbPacket)*100)
}