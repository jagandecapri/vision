package preprocess

import (
	"github.com/google/gopacket"
	"github.com/jagandecapri/vision/tree"
	"github.com/jagandecapri/vision/preprocess/aggregates"
	"log"
)

type MicroSlot []tree.Point

type X_micro_slot struct{
	AggSrc MicroSlot
	AggDst MicroSlot
	AggSrcDst MicroSlot
}

type Accumulator struct{
	AggSrc map[gopacket.Endpoint]aggregates.AggSrc
	AggDst map[gopacket.Endpoint]aggregates.AggDst
	AggSrcDst map[gopacket.Flow]aggregates.AggSrcDst
}

func (acc *Accumulator) AddPacket(p gopacket.Packet){
	networkLayer := p.NetworkLayer()

	if networkLayer == nil{
		log.Println("Packet has no network layer")
	} else {
		netFlow := networkLayer.NetworkFlow()

		src := netFlow.Src()
		dst := netFlow.Dst()
		srcdst := netFlow

		aggsrc, ok := acc.AggSrc[src]
		if !ok{
			aggsrc = aggregates.NewAggSrc(src)
		}
		aggsrc.AddPacket(p)
		acc.AggSrc[src] = aggsrc

		aggdst, ok := acc.AggDst[dst]
		if !ok{
			aggdst = aggregates.NewAggDst(dst)
		}
		aggdst.AddPacket(p)
		acc.AggDst[dst] = aggdst

		aggsrcdst, ok := acc.AggSrcDst[netFlow]
		if !ok{
			aggsrcdst = aggregates.NewAggSrcDst(src, dst)
		}
		aggsrcdst.AddPacket(p)
		acc.AggSrcDst[srcdst] = aggsrcdst
	}
}

func (acc *Accumulator) extractFeatures(aggInterface aggregates.AggInterface) (tree.Point, error){
	x := map[string]float64{
		"nbPacket": aggInterface.NbPacket(),
		"nbSrcPort": aggInterface.NbSrcPort(),
		"nbDstPort": aggInterface.NbDstPort(),
		"nbSrcs": aggInterface.NbSrcs(),
		"nbDsts": aggInterface.NbDsts(),
		"perSYN": aggInterface.PerSYN(),
		"perACK": aggInterface.PerACK(),
		"perICMP": aggInterface.PerICMP(),
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
			//log.Println(x_aggdst)
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
	return Accumulator{AggSrc: make(map[gopacket.Endpoint]aggregates.AggSrc),
		AggDst: make(map[gopacket.Endpoint]aggregates.AggDst),
		AggSrcDst: make(map[gopacket.Flow]aggregates.AggSrcDst),
	}
}
