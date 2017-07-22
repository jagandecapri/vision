package tree

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

func (p *PointContainer) Distance(point Point) float64{
	return 0.0
}

func (p *PointContainer) PlaneDistance(val float64, dim int) float64{
	return 0.0
}
