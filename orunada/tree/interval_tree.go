package tree

import (
	"github.com/golang-collections/go-datastructures/augmentedtree"
)

type IntervalConc struct {
	Id   int
	Low  []int64
	High []int64
}

func (itv IntervalConc) LowAtDimension(dim uint64) int64{
	return itv.Low[dim - 1]
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (itv IntervalConc) HighAtDimension(dim uint64) int64{
	return itv.High[dim - 1]
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (itv IntervalConc) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool{
	check := false
	for i := uint64(1); i <= uint64(len(itv.Low)); i++{
		if interval.LowAtDimension(i) <= itv.HighAtDimension(i) &&
			interval.HighAtDimension(i) >= itv.LowAtDimension(i){
			check = true
		} else {
			check = false
		}
	}
	return check
}

// ID should be a unique ID representing this interval.  This
// is used to identify which interval to delete from the tree if
// there are duplicates.
func (itv IntervalConc) ID() uint64{
	return uint64(itv.Id)
}

// Return the total number of dimensions
func (itv IntervalConc) Dim() int{
	return len(itv.Low)
}

// Return the value X_{dim}, dim is started from 0
func (itv IntervalConc) GetValue(dim int) int{
	dim++
	res_int64 := itv.HighAtDimension(uint64(dim))
	return int(res_int64)
}

// Return the distance between two points
func (itv IntervalConc) Distance(point Point) float64{
	return 0.0
}

// Return the distance between the point and the plane X_{dim}=val
func (itv IntervalConc) PlaneDistance(val float64, dim int) float64{
	return 0.0
}

func IntervalBuilder(min int, max int, interval_length int) []IntervalConc {
	id := 1
	intervals := []IntervalConc{}
	for i := min; i < max; i += interval_length{
		for j := min; j < max; j += interval_length{
			intervals = append(intervals, IntervalConc{	Id: id,
				Low: []int64{int64(i), int64(j)},
				High: []int64{int64(i + interval_length), int64(j + interval_length)}})
			id += 1
		}
	}
	return intervals
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}