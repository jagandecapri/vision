package tree

import "github.com/golang-collections/go-datastructures/rangetree"

type Entry struct{
	point []int64
}

func (e Entry) ValueAtDimension(dimension uint64) int64{
	return e.point[dimension]
}

type Interval struct {
	low []int64
	high []int64
}

// LowAtDimension returns an integer representing the lower bound
// at the requested dimension.
func (i Interval) LowAtDimension(dimension uint64) int64{
	return i.low[dimension]
}

func (i Interval) HighAtDimension(dimension uint64) int64{
	return i.high[dimension]
}

func NewRangeTree(dim uint64) rangetree.RangeTree{
	return rangetree.New(dim)
}