package tree

import "github.com/golang-collections/go-datastructures/augmentedtree"

type IntervalConc struct {
	id int
	low []int64
	high []int64
}

func (itv IntervalConc) LowAtDimension(dim uint64) int64{
	return itv.low[dim - 1]
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (itv IntervalConc) HighAtDimension(dim uint64) int64{
	return itv.high[dim - 1]
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (itv IntervalConc) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool{
	interval = interval.(IntervalConc)
	if itv.LowAtDimension(dim) <= interval.LowAtDimension(dim) &&
		itv.HighAtDimension(dim) > interval.HighAtDimension(dim){
		return true
	} else {
		return false
	}
}
// ID should be a unique ID representing this interval.  This
// is used to identify which interval to delete from the tree if
// there are duplicates.
func (itv IntervalConc) ID() uint64{
	return uint64(itv.id)
}

func IntervalBuilder(min int, max int, interval_length int) []IntervalConc {
	id := 0
	intervals := []IntervalConc{}
	for i := min; i < max; i += interval_length{
		for j := min; j < max; j += interval_length{
			intervals = append(intervals, IntervalConc{	id: id,
				low: []int64{int64(i), int64(i + interval_length)},
				high: []int64{int64(j), int64(j + interval_length)}})
			id += 1
		}
	}
	return intervals
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}