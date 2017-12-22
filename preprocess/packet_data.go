package preprocess

import "github.com/google/gopacket"

type PacketData struct{
	Data     gopacket.Packet
	Metadata gopacket.CaptureInfo
}
