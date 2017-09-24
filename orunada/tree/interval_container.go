package tree

import (
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"math"
)

type Range struct{
	Low          [2]float64
	High         [2]float64
}

type IntervalContainer struct {
	Id           int
	Scale_factor int
	*Range
}

func (itv IntervalContainer) LowAtDimension(dim uint64) int64{
	tmp :=  int64(itv.Low[dim - 1] * math.Pow(10, float64(itv.Scale_factor)))
	return tmp
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (itv IntervalContainer) HighAtDimension(dim uint64) int64{
	tmp := int64(itv.High[dim - 1] * math.Pow(10, float64(itv.Scale_factor)))
	return tmp
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (itv IntervalContainer) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool{
	check := false
	//for i := uint64(1); i <= uint64(len(itv.Low)); i++{
	//	if interval.LowAtDimension(i) <= itv.HighAtDimension(i) &&
	//		interval.HighAtDimension(i) >= itv.LowAtDimension(i){
	//		check = true
	//	} else {
	//		check = false
	//	}
	//}
	//for i := uint64(1); i <= uint64(len(itv.Low)); i++{
		if interval.LowAtDimension(dim) <= itv.HighAtDimension(dim) &&
			interval.HighAtDimension(dim) >= itv.LowAtDimension(dim){
			check = true
		} else {
			check = false
		}
	//}
	return check
}

// ID should be a unique ID representing this interval.  This
// is used to identify which interval to delete from the tree if
// there are duplicates.
func (itv IntervalContainer) ID() uint64{
	return uint64(itv.Id)
}

func (itv IntervalContainer) GetCenter() []float64{
	dim := len(itv.Low)
	center := []float64{}
	for i := 0; i < dim; i++{
		center = append(center, (itv.High[i] - itv.Low[i])/2 )
	}
	return center
}