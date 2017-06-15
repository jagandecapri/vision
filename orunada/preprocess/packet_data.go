package preprocess

import "github.com/google/gopacket"

type PacketData struct{
	data gopacket.Packet
	metadata gopacket.CaptureInfo
}
