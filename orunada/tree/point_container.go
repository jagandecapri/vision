package tree

import "math"

type PointContainer struct{
	dim int
	point []int
}

func (p PointContainer) Dim() int{
	return p.dim
}

func (p PointContainer) GetValue(dim int) int{
	return p.point[dim]
}

func (p *PointContainer) Distance(p1 *PointContainer) float64{
	sum := 0.0
	for i:=0; i<len(p.point); i++{
		sum += math.Pow(float64(p.point[i]-p1.point[i]), 2)
	}
	euclidean_dist := math.Sqrt(sum)
	return euclidean_dist
}

func (p *PointContainer) PlaneDistance(val float64, dim int) float64{
	return 0.0
}
