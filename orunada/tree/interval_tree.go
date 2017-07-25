package tree

import "github.com/Workiva/go-datastructures/augmentedtree"

type Interval struct {
	id int
	low []int64
	high []int64
}

// LowAtDimension returns an integer representing the lower bound
// at the requested dimension.
func (itv *Interval) LowAtDimension(dim uint64) int64{
	return itv.low[dim - 1]
}

// HighAtDimension returns an integer representing the higher bound
// at the requested dimension.
func (itv *Interval) HighAtDimension(dim uint64) int64{
	return itv.high[dim - 1]
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (itv *Interval) OverlapsAtDimension(interval Interval, dim uint64) bool{
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
func (itv *Interval) ID() uint64{
	return uint64(itv.id)
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}