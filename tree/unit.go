package tree

import (
	"math"
)

type Unit struct {
	Id                 int
	Cluster_id         int
	Dimension          int
	Center             PointContainer
	Neighbour_units    map[Range]*Unit
	points             map[int]PointContainer
	Center_calculated bool
	Range
}

func (u *Unit) AddPoint(p PointContainer) {
	u.Center_calculated = false
	u.points[p.GetID()] = p
}

func (u *Unit) RemovePoint(p PointContainer) {
	u.Center_calculated = false
	delete(u.points, p.GetID())
}

func (u *Unit) CalculateCenter() {
	u.GetCenter()
}

func (u *Unit) GetCenter() PointContainer {
	if len(u.Center.Vec) > 0{
		return u.Center
	}
	Center_vec := make([]float64, u.Dimension)
	for _, p := range u.points {
		for i := 0; i < p.Dim(); i++ {
			Center_vec[i] = Center_vec[i] + p.GetValue(i)
		}
	}
	for i, _ := range Center_vec {
		Center_vec[i] = Center_vec[i] / float64(len(u.points))
	}
	u.Center = PointContainer{Unit_id: u.Id, Vec: Center_vec}
	return u.Center
}

func (u *Unit) GetNumberOfPoints() int {
	return len(u.points)
}

//Implementing PointInterface methods
func (u *Unit) GetID() int {
	return u.Id
}

func (u *Unit) Dim() int {
	return len(u.Center.Vec)
}

func (u *Unit) GetValue(dim int) float64 {
	return u.Center.Vec[dim]
}

func (u *Unit) Distance(p1 PointInterface) float64 {
	if u.Center_calculated {
		u.CalculateCenter()
	}
	sum := 0.0
	t := p1.(*PointContainer)
	for i := 0; i < len(u.Center.Vec); i++ {
		sum += math.Pow(u.Center.Vec[i]-t.Vec[i], 2)
	}
	euclidean_dist := math.Sqrt(sum)
	return euclidean_dist
}

func (u *Unit) PlaneDistance(val float64, dim int) float64 {
	return 0.0
}

func NewUnit(id int, dimension int, rg Range) Unit{
	unit := Unit{
		Id: id,
		Dimension: dimension,
		Neighbour_units: make(map[Range]*Unit),
		points: make(map[int]PointContainer),
		Range: rg,
	}
	return unit
}