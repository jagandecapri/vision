package preprocess

import (
	"testing"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"github.com/stretchr/testify/assert"
)

func TestAccumulator_GetMicroSlot(t *testing.T) {
	opts := gopacket.SerializeOptions{FixLengths: true}

	buf1 := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf1, opts,
		&layers.Ethernet{
			SrcMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthernetType: layers.EthernetTypeIPv4},
		&layers.IPv4{  SrcIP: net.IP{1, 2, 3, 4},
			DstIP: net.IP{5, 6, 7, 8},
			Protocol: layers.IPProtocolICMPv4},
		&layers.ICMPv4{TypeCode: layers.ICMPv4TypeRedirect},
		gopacket.Payload([]byte{1, 2, 3, 4}))
	packetData1 := buf1.Bytes()
	packet1 := gopacket.NewPacket(packetData1, layers.LayerTypeEthernet, gopacket.Default)

	buf2 := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf2, opts,
		&layers.Ethernet{
			SrcMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthernetType: layers.EthernetTypeIPv4},
		&layers.IPv4{  SrcIP: net.IP{1, 2, 3, 4},
			DstIP: net.IP{5, 6, 7, 8},
			Protocol: layers.IPProtocolTCP},
		&layers.TCP{SrcPort: 1234, DstPort: 5678, ACK: true},
		gopacket.Payload([]byte{1, 2, 3, 4}))
	packetData2 := buf2.Bytes()
	packet2 := gopacket.NewPacket(packetData2, layers.LayerTypeEthernet, gopacket.Default)

	buf3 := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf3, opts,
		&layers.Ethernet{
			SrcMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthernetType: layers.EthernetTypeIPv4},
		&layers.IPv4{  SrcIP: net.IP{1, 2, 3, 4},
			DstIP: net.IP{5, 6, 7, 8},
			Protocol: layers.IPProtocolTCP},
		&layers.TCP{SrcPort: 1234, DstPort: 5678, ACK: true},
		gopacket.Payload([]byte{1, 2, 3, 4}))
	packetData3 := buf3.Bytes()
	packet3 := gopacket.NewPacket(packetData3, layers.LayerTypeEthernet, gopacket.Default)

	acc := NewAccumulator()
	acc.AddPacket(packet1)
	acc.AddPacket(packet2)
	acc.AddPacket(packet3)

	output := acc.AggDst[packet1.NetworkLayer().NetworkFlow().Dst()]
	assert.Equal(t, 3.0, output.NbPacket())
	assert.InDelta(t, 66.666, output.PerACK(), 66.1)
	assert.InDelta(t, 33.333, output.PerICMP(), 33.1)
}
