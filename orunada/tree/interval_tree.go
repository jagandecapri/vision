package tree

import (
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"math"
)

type IntervalConc struct {
	Id   int
	Interval_conc_store []IntervalConc
	Low  []float64
	High []float64
	Decimal_places int
}

func (itv IntervalConc) LowAtDimension(dim uint64) int64{
	return int64(itv.Low[dim - 1] * math.Pow(10, float64(itv.Decimal_places)))
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (itv IntervalConc) HighAtDimension(dim uint64) int64{
	return int64(itv.High[dim - 1] * math.Pow(10, float64(itv.Decimal_places)))
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

func IntervalBuilder(min int, max int, interval_length int) []IntervalConc {
	id := 1
	intervals := []IntervalConc{}
	for i := min; i < max; i += interval_length{
		for j := min; j < max; j += interval_length{
			intervals = append(intervals, IntervalConc{	Id: id,
				Low: []float64{float64(i), float64(j)},
				High: []float64{float64(i + interval_length), float64(j + interval_length)}})
			id += 1
		}
	}
	return intervals
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}