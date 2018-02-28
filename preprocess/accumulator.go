package preprocess

import (
	"github.com/google/gopacket"
	"github.com/jagandecapri/vision/tree"
	"errors"
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

func (acc *Accumulator) extractFeatures(aggInterface AggInterface) (tree.Point, tree.Point, tree.Point, error){
	if aggInterface.NbPacket() == 1{
		return tree.Point{}, tree.Point{}, tree.Point{},  errors.New("NbPacket is 1, normal packet") //Assumption taken here
	}

	x_aggsrc := map[string]float64{
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
	pnt_aggsrc := tree.Point{Id: point_ctr, Vec_map: x_aggsrc}

	x_aggdst := map[string]float64{
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
	pnt_aggdst := tree.Point{Id: point_ctr, Vec_map: x_aggdst}

	x_aggsrcdst:= map[string]float64{
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
	pnt_aggsrcdst := tree.Point{Id: point_ctr, Vec_map: x_aggsrcdst}

	return pnt_aggsrc, pnt_aggdst, pnt_aggsrcdst, nil
}

func (acc *Accumulator) GetMicroSlot() X_micro_slot{
	X := X_micro_slot{}

	for _, val := range acc.AggSrc{
		x_aggsrc, x_aggdst, x_aggsrcdst, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggSrc = append(X.AggSrc, x_aggsrc)
			X.AggDst = append(X.AggDst, x_aggdst)
			X.AggSrcDst = append(X.AggSrcDst, x_aggsrcdst)
		}
	}

	for _, val := range acc.AggDst{
		x_aggsrc, x_aggdst, x_aggsrcdst, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggSrc = append(X.AggSrc, x_aggsrc)
			X.AggDst = append(X.AggDst, x_aggdst)
			X.AggSrcDst = append(X.AggSrcDst, x_aggsrcdst)
		}
	}

	for _, val := range acc.AggSrcDst{
		x_aggsrc, x_aggdst, x_aggsrcdst, err := acc.extractFeatures(&val)
		if err == nil{
			X.AggSrc = append(X.AggSrc, x_aggsrc)
			X.AggDst = append(X.AggDst, x_aggdst)
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
	//return Accumulator{AggSrc: make(map[string]AggSrc),
	//	AggDst: make(map[string]AggDst),
	//	AggSrcDst: make(map[string]AggSrcDst),
	//}
}
