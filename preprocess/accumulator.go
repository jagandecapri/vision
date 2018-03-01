package preprocess

import (
	"github.com/google/gopacket"
	"github.com/jagandecapri/vision/tree"
)

type X_micro_slot struct{
	AggSrc []tree.Point
	AggDst []tree.Point
	AggSrcDst []tree.Point
}

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

func (acc *Accumulator) extractFeatures(aggInterface AggInterface) (tree.Point, error){
	//if aggInterface.NbPacket() == 1{
	//	return tree.Point{}, errors.New("NbPacket is 1, normal packet") //Assumption taken here
	//}

	x := map[string]float64{
		"nbPacket": aggInterface.NbPacket(),
		"nbSrcPort": aggInterface.NbSrcPort(),
		"nbDstPort": aggInterface.NbDstPort(),
		"nbSrcs": aggInterface.NbSrcs(),
		"nbDsts": aggInterface.NbDsts(),
		"perSyn": aggInterface.PerSYN(),
		"perAck": aggInterface.PerACK(),
		"perRST": aggInterface.PerRST(),
		"perFIN": aggInterface.PerFIN(),
		"perCWR": aggInterface.PerCWR(),
		"perURG": aggInterface.PerURG(),
		"avgPktSize": aggInterface.AvgPktSize(),
		"meanTTL": aggInterface.MeanTTL(),
	}

	Point_ctr = Point_ctr + 1
	pnt := tree.Point{Id: Point_ctr, Key: aggInterface.GetKey(), Vec_map: x}

	return pnt, nil
}

func (acc *Accumulator) GetMicroSlot() X_micro_slot{
	X := X_micro_slot{}

	for _, val := range acc.AggSrc{
		x_aggsrc, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggSrc = append(X.AggSrc, x_aggsrc)
		}
	}

	for _, val := range acc.AggDst{
		x_aggdst, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggDst = append(X.AggDst, x_aggdst)
		}
	}

	for _, val := range acc.AggSrcDst{
		x_aggsrcdst, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggSrcDst = append(X.AggSrcDst, x_aggsrcdst)
		}
	}

	return X
}

func NewAccumulator() Accumulator{
	return Accumulator{AggSrc: make(map[gopacket.Endpoint]AggSrc),
		AggDst: make(map[gopacket.Endpoint]AggDst),
		AggSrcDst: make(map[gopacket.Flow]AggSrcDst),
	}
}
