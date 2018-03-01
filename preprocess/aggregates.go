package preprocess

import (
	"github.com/google/gopacket/layers"
	"net"
	"github.com/google/gopacket"
)

type FlowKey struct{
	SrcIP, DstIP net.IP
	SrcPort, DstPort layers.TCPPort
}

type AggInterface interface{
	GetKey() []FlowKey
	NbPacket() float64
	NbSrcPort() float64
	NbDstPort() float64
	NbSrcs() float64
	NbDsts() float64
	PerSYN() float64
	PerACK() float64
	PerRST() float64
	PerFIN() float64
	PerCWR() float64
	PerURG() float64
	AvgPktSize() float64
	MeanTTL() float64
	SimIPDst() float64
	SimIPSrc() float64
	PerICMP() float64
	PerICMPRed() float64
	PerICMPTime() float64
	PerICMPUnr() float64
	PerICMPOther() float64
}

func NewAggSrc() AggSrc{
	aggSrc := AggSrc{dsts: make(map[gopacket.Endpoint]int),
		dsts_subnetwork: make(map[gopacket.Endpoint][]int)}
	return aggSrc
}

type AggSrc struct{
	FlowKeys []FlowKey
	dsts map[gopacket.Endpoint]int
	dsts_subnetwork map[gopacket.Endpoint][] int
}

func (a *AggSrc) AddPacket(p gopacket.Packet) gopacket.ErrorLayer{
	dst := 	p.NetworkLayer().NetworkFlow().Dst()
	a.dsts[dst] = 1

	var srcIP, dstIP net.IP
	var srcPort, dstPort layers.TCPPort

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		srcIP = ipv4.SrcIP
		dstIP = ipv4.DstIP
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		srcIP = ipv6.SrcIP
		dstIP = ipv6.DstIP
	}

	// Let's see if the packet is TCP
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = tcp.SrcPort
		dstPort = tcp.DstPort
	}

	flowKey := FlowKey{SrcIP: srcIP, DstIP: dstIP, SrcPort: srcPort, DstPort: dstPort}
	a.FlowKeys = append(a.FlowKeys, flowKey)

	// Check for errors
	var err gopacket.ErrorLayer
	if err = p.ErrorLayer(); err != nil {
		return err
	}
	return err
}

func (a *AggSrc) GetKey() []FlowKey{
	return a.FlowKeys
}

func (a *AggSrc) NbPacket() float64 {
	return float64(len(a.dsts))
}

func (a *AggSrc) NbSrcPort() float64 {
	return 0.0
}

func (a *AggSrc) NbDstPort() float64 {
	return float64(len(a.dsts))
}

func (a *AggSrc) NbSrcs() float64 {
	return 0.0
}

func (a *AggSrc) NbDsts() float64 {
	return 0.0
}

func (a *AggSrc) PerSYN() float64 {
	return 0.0
}

func (AggSrc) PerACK() float64 {
	return 0.0
}

func (AggSrc) PerRST() float64 {
	return 0.0
}

func (AggSrc) PerFIN() float64 {
	return 0.0
}

func (AggSrc) PerCWR() float64 {
	return 0.0
}

func (AggSrc) PerURG() float64 {
	return 0.0
}

func (AggSrc) AvgPktSize() float64 {
	return 0.0
}

func (AggSrc) MeanTTL() float64 {
	return 0.0
}

func (AggSrc) SimIPDst() float64 {
	return 0.0
}

func (AggSrc) SimIPSrc() float64 {
	return 0.0
}

func (AggSrc) PerICMP() float64{
	return 0.0
}

func (AggSrc) PerICMPRed() float64 {
	return 0.0
}

func (AggSrc) PerICMPTime() float64 {
	return 0.0
}

func (AggSrc) PerICMPUnr() float64 {
	return 0.0
}

func (AggSrc) PerICMPOther() float64 {
	return 0.0
}

func NewAggDst() AggDst{
	aggDst := AggDst{srcs: make(map[gopacket.Endpoint]int)}
	return aggDst
}

type AggDst struct{
	FlowKeys []FlowKey
	srcs map[gopacket.Endpoint]int
	srcPort map[layers.TCPPort]int
	nbPacket float64
	nbSYN float64
	nbICMP float64
}

func (a *AggDst) AddPacket(p gopacket.Packet) gopacket.ErrorLayer{

	a.nbPacket++

	src := 	p.NetworkLayer().NetworkFlow().Src()
	a.srcs[src] = 1

	var srcIP, dstIP net.IP
	var srcPort, dstPort layers.TCPPort

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		srcIP = ipv4.SrcIP
		dstIP = ipv4.DstIP
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		srcIP = ipv6.SrcIP
		dstIP = ipv6.DstIP
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
	}

	flowKey := FlowKey{SrcIP: srcIP, DstIP: dstIP, SrcPort: srcPort, DstPort: dstPort}
	a.FlowKeys = append(a.FlowKeys, flowKey)

	// Let's see if the packet is icmp4
	icmp4Layer :=p.Layer(layers.LayerTypeICMPv4)
	if icmp4Layer != nil {
		a.nbICMP++
	}

	// Let's see if the packet is icmp6
	icmp6Layer := p.Layer(layers.LayerTypeICMPv6)
	if icmp6Layer != nil {
		a.nbICMP++
	}

	// Check for errors
	var err gopacket.ErrorLayer
	if err = p.ErrorLayer(); err != nil {
		return err
	}
	return err
}

func (a *AggDst) GetKey() []FlowKey{
	return a.FlowKeys
}

func (a *AggDst) NbPacket() float64 {
	return float64(len(a.srcs))
}

func (a *AggDst) NbSrcPort() float64 {
	return 0.0
}

func (a *AggDst) NbDstPort() float64 {
	return 0.0
}

func (a *AggDst) NbSrcs() float64 {
	return float64(len(a.srcs))
}

func (a *AggDst) NbDsts() float64 {
	return 0.0
}

func (a *AggDst) PerSYN() float64 {
	return float64((a.nbICMP/a.nbPacket)*100)
}

func (a *AggDst) PerACK() float64 {
	return 0.0
}

func (a *AggDst) PerRST() float64 {
	return 0.0
}

func (a *AggDst) PerFIN() float64 {
	return 0.0
}

func (a *AggDst) PerCWR() float64 {
	return 0.0
}

func (a *AggDst) PerURG() float64 {
	return 0.0
}

func (a *AggDst) AvgPktSize() float64 {
	return 0.0
}

func (a *AggDst) MeanTTL() float64 {
	return 0.0
}

func (a *AggDst) SimIPDst() float64 {
	return 0.0
}

func (a *AggDst) SimIPSrc() float64 {
	return 0.0
}

func (a *AggDst) PerICMP() float64{
	return float64((a.nbICMP/a.nbPacket)*100)
}

func (a *AggDst) PerICMPRed() float64 {
	return 0.0
}

func (a *AggDst) PerICMPTime() float64 {
	return 0.0
}

func (a *AggDst) PerICMPUnr() float64 {
	return 0.0
}

func (a *AggDst) PerICMPOther() float64 {
	return 0.0
}

func NewAggSrcDst() AggSrcDst {
	aggDstSrc := AggSrcDst{srcPort: make(map[layers.TCPPort]int),
		dstPort: make(map[layers.TCPPort]int)}
	return aggDstSrc
}

type AggSrcDst struct{
	FlowKeys []FlowKey
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

	flowKey := FlowKey{SrcIP: srcIP, DstIP: dstIP, SrcPort: srcPort, DstPort: dstPort}
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

func (a *AggSrcDst) GetKey() []FlowKey{
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