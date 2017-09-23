package tree

import (
	"math"
)

type Unit struct{
	Id int
	Dimension int
	Center PointContainer
	points map[int]PointContainer
	recalculate_Center bool
	*IntervalContainer
}

func (u *Unit) AddPoint(p PointContainer){
	u.recalculate_Center = true
	u.points[p.GetID()] = p
}

func (u *Unit) RemovePoint(){

}

func (u *Unit) RecalculateCenter(){
	u.GetCenter()
}

func (u *Unit) GetCenter() PointContainer{
	Center_vec := make([]float64, u.Dimension)
	for _, p := range u.points{
		for i := 0; i < p.Dim(); i++{
			Center_vec[i] = Center_vec[i] + p.GetValue(i)
		}
	}
	for i, _ := range Center_vec{
		Center_vec[i] = Center_vec[i]/float64(len(u.points))
	}
	u.Center = PointContainer{Unit_id: u.Id, Vec: Center_vec}
	return u.Center
}

func (u *Unit) GetNumberOfPoints() int{
	return len(u.points)
}


//Implementing PointInterface methods
func (u *Unit) GetID() int{
	return u.Id
}

func (u *Unit) Dim() int{
	return len(u.Center.Vec)
}

func (u *Unit) GetValue(dim int) float64{
	return u.Center.Vec[dim]
}

func (u *Unit) Distance(p1 PointInterface) float64{
	if u.recalculate_Center{
		u.RecalculateCenter()
	}
	sum := 0.0
	t := p1.(*PointContainer)
	for i:=0; i<len(u.Center.Vec); i++{
		sum += math.Pow(u.Center.Vec[i]-t.Vec[i], 2)
	}
	euclidean_dist := math.Sqrt(sum)
	return euclidean_dist
}

func (u *Unit) PlaneDistance(val float64, dim int) float64{
	return 0.0
}


//func (u *Unit) LowAtDimension(dim uint64) int64{
//	return int64(u.Low[dim - 1] * math.Pow(10, float64(u.Decimal_places)))
//}
//
//// HighAtDimension returns an integer representing the higher bound
//// at the requested dimension.
//func (u *Unit) HighAtDimension(dim uint64) int64{
//	return int64(u.High[dim - 1] * math.Pow(10, float64(u.Decimal_places)))
//}
//
//// OverlapsAtDimension should return a bool indicating if the provided
//// interval overlaps this interval at the dimension requested.
//func (u *Unit) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool{
//	check := false
//	for i := uint64(1); i <= uint64(len(u.Low)); i++{
//		if interval.LowAtDimension(i) <= u.HighAtDimension(i) &&
//			interval.HighAtDimension(i) >= u.LowAtDimension(i){
//			check = true
//		} else {
//			check = false
//		}
//	}
//	return check
//}
//
//// ID should be a unique ID representing this interval.  This
//// is used to identify which interval to delete from the tree if
//// there are duplicates.
//func (u *Unit) ID() uint64{
//	return uint64(u.Id)
//}