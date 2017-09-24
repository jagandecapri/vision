package tree

import "math"

type PointContainer struct{
	Unit_id int
	Vec []float64
	Point
}

func (p *PointContainer) GetID() int{
	return p.Id
}

func (p *PointContainer) Dim() int{
	return len(p.Vec)
}

func (p *PointContainer) GetValue(dim int) float64{
	return p.Vec[dim]
}

func (p *PointContainer) Distance(p1 PointInterface) float64{
	sum := 0.0
	t := p1.(*PointContainer)
	for i:=0; i<len(p.Vec); i++{
		sum += math.Pow(p.Vec[i]-t.Vec[i], 2)
	}
	euclidean_dist := math.Sqrt(sum)
	return euclidean_dist
}

func (p *PointContainer) PlaneDistance(val float64, dim int) float64{
	return 0.0
}
