package preprocess

import (
	"github.com/google/gopacket"
)


type Accumulator struct{
	AggSrc map[gopacket.Endpoint]AggSrc
	AggDst map[gopacket.Endpoint]AggDst
	AggSrcDst map[gopacket.Flow]AggSrcDst
}

func (acc *Accumulator) AddPacket(p gopacket.Packet){
	netFlow := p.NetworkLayer().NetworkFlow()

	src := netFlow.Src()
	dst := netFlow.Dst()
	srcdst := netFlow

	aggsrc, ok := acc.AggSrc[src]
	if !ok{
		aggsrc = NewAggSrc()
	}
	aggsrc.AddPacket(p)
	acc.AggSrc[src] = aggsrc

	aggdst, ok := acc.AggDst[dst]
	if !ok{
		aggdst = NewAggDst()
	}
	aggdst.AddPacket(p)
	acc.AggDst[dst] = aggdst

	aggsrcdst, ok := acc.AggSrcDst[netFlow]
	if !ok{
		aggsrcdst = NewAggSrcDst()
	}
	aggsrcdst.AddPacket(p)
	acc.AggSrcDst[srcdst] = aggsrcdst
}

func NewAccumulator() Accumulator{
	return Accumulator{AggSrc: make(map[gopacket.Endpoint]AggSrc),
		AggDst: make(map[gopacket.Endpoint]AggDst),
		AggSrcDst: make(map[gopacket.Flow]AggSrcDst),
	}
}
