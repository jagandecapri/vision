package preprocess

import (
	"net"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
)

type PacketFeature struct{
	SrcIP, DstIP net.IP
	SrcPort, DstPort layers.TCPPort
	SYN, ACK, RST, FIN, CWR, URG bool
	length uint16
	TTL uint8
}

func (p *PacketFeature) ExtractFeature(rawPacket gopacket.Packet) *PacketFeature{
	// Let's see if the packet is an ethernet packet
	ethernetLayer := rawPacket.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		// ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
	}

	// Let's see if the packet is IPv4 (even though the ether type told us)
	ipv4Layer := rawPacket.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		p.SrcIP = ipv4.SrcIP
		p.DstIP = ipv4.DstIP
		p.TTL = ipv4.TTL
	}
	// Let's see if the packet is IPv6 (even though the ether type told us)
	ipv6Layer := rawPacket.Layer(layers.LayerTypeIPv4)
	if ipv6Layer != nil {
		//fmt.Println("IPv4 layer detected.")
		ipv6, _ := ipv6Layer.(*layers.IPv4)
		p.SrcIP = ipv6.SrcIP
		p.DstIP = ipv6.DstIP
		p.TTL = ipv6.TTL
	}

	// Let's see if the packet is TCP
	tcpLayer := rawPacket.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		p.SrcPort = tcp.SrcPort
		p.DstPort = tcp.DstPort
		p.SYN = tcp.SYN
		p.ACK = tcp.ACK
		p.RST = tcp.RST
		p.FIN = tcp.FIN
		p.CWR = tcp.CWR
		p.URG = tcp.URG
	}

	// Let's see if the packet is icmp4
	icmp4Layer := rawPacket.Layer(layers.LayerTypeICMPv4)
	if icmp4Layer != nil {
		//fmt.Println("ICMP4 Layer Detected")
		//icmp4, _ := icmp4Layer.(*layers.ICMPv4)
		//fmt.Println(icmp4.TypeCode.String())
	}

	// Let's see if the packet is icmp6
	icmp6Layer := rawPacket.Layer(layers.LayerTypeICMPv6)
	if icmp6Layer != nil {
		//fmt.Println("ICMP6 Layer Detected")
		//icmp6, _ := icmp6Layer.(*layers.ICMPv6)
		//fmt.Println(icmp6.TypeCode.String())
	}

	// Check for errors
	if err := rawPacket.ErrorLayer(); err != nil {
		//fmt.Println("Error decoding some part of the packet:", err)
	}

	return p
}