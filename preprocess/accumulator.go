package preprocess

import (
	"github.com/google/gopacket"
	"github.com/jagandecapri/vision/tree"
)

type X_micro_slot []tree.Point

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

func (acc *Accumulator) extractFeatures(aggInterface AggInterface) tree.Point{
	x := map[string]float64{
		"nbPacket": aggInterface.NbPacket(),
		"nbSrcPort": aggInterface.NbSrcPort(),
		"nbDstPort": aggInterface.NbDstPort(),
		"nbSrcs": aggInterface.NbSrcs(),
		"nbDsts": aggInterface.NbDsts(),
		"perSyn": aggInterface.PerSyn(),
		"perAck": aggInterface.PerAck(),
		"perRST": aggInterface.PerRST(),
		"perFIN": aggInterface.PerFIN(),
		"perCWR": aggInterface.PerCWR(),
		"perURG": aggInterface.PerURG(),
		"avgPktSize": aggInterface.AvgPktSize(),
		"meanTTL": aggInterface.MeanTTL(),
	}
	point_ctr = point_ctr + 1
	pnt := tree.Point{Id: point_ctr, Vec_map: x}
	return pnt
}

func (acc *Accumulator) GetMicroSlot() X_micro_slot{
	X := X_micro_slot{{}}

	for _, val := range acc.AggSrc{
		x := acc.extractFeatures(&val)
		X = append(X, x)
	}

	for _, val := range acc.AggDst{
		x := acc.extractFeatures(&val)
		X = append(X, x)
	}

	for _, val := range acc.AggDst{
		x := acc.extractFeatures(&val)
		X = append(X, x)
	}

	return X
}

func NewAccumulator() Accumulator{
	return Accumulator{AggSrc: make(map[gopacket.Endpoint]AggSrc),
		AggDst: make(map[gopacket.Endpoint]AggDst),
		AggSrcDst: make(map[gopacket.Flow]AggSrcDst),
	}
}
