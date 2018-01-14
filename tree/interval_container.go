package tree

import (
	"github.com/Workiva/go-datastructures/augmentedtree"
	"math"
)

type IntervalContainer struct {
	Id           int
	Scale_factor int
	Range
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
	if interval.LowAtDimension(dim) < itv.HighAtDimension(dim) &&
		interval.HighAtDimension(dim) >= itv.LowAtDimension(dim) {
		return true
	}

	return false
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