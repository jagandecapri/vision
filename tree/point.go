package tree

import (
	"math"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
)

type PointKey struct{
	SrcIP, DstIP []gopacket.Endpoint
	SrcPort, DstPort []layers.TCPPort
}

type Point struct{
	Id       int
	Key PointKey
	Unit_id int
	Vec []float64
	Vec_map  map[string]float64
}

func (p *Point) GetID() int{
	return p.Id
}

func (p *Point) Dim() int{
	return len(p.Vec)
}

func (p *Point) GetValue(dim int) float64{
	return p.Vec[dim]
}

func (p *Point) Distance(p1 PointInterface) float64{
	sum := 0.0
	t := p1.(*Point)
	for i:=0; i<len(p.Vec); i++{
		sum += math.Pow(p.Vec[i]-t.Vec[i], 2)
	}
	euclidean_dist := math.Sqrt(sum)
	return euclidean_dist
}

func (p *Point) PlaneDistance(val float64, dim int) float64{
	return 0.0
}
