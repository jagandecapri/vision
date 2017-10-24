package tree

import (
	"testing"
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/stretchr/testify/assert"
)

type IntervalContainerTest struct {
	Id   int
	Low  []int64
	High []int64
}

func (itv IntervalContainerTest) LowAtDimension(dim uint64) int64 {
	return itv.Low[dim-1]
}

func (itv IntervalContainerTest) HighAtDimension(dim uint64) int64 {
	return itv.High[dim-1]
}

// OverlapsAtDimension should return a bool indicating if the provided
// interval overlaps this interval at the dimension requested.
func (itv IntervalContainerTest) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool {
	if interval.LowAtDimension(dim) < itv.HighAtDimension(dim) &&
		interval.HighAtDimension(dim) >= itv.LowAtDimension(dim) {
		return true
	}

	return false
}

func (itv IntervalContainerTest) ID() uint64 {
	return uint64(itv.Id)
}

func TestIntervalQueryInt(t *testing.T) {
	tree := augmentedtree.New(2)
	intervals_container := IntervalContainerTest{Id: 1, Low: []int64{0, 0}, High: []int64{1, 1}}
	intervals := augmentedtree.Interval(intervals_container)
	tree.Add(intervals)
	interval_test_query := tree.Query(IntervalContainerTest{Id: 1, Low: []int64{0, 0}, High: []int64{0, 0}})
	assert.Len(t, interval_test_query, 1, "Interval not found: %v", interval_test_query)
	interval_test_query1 := tree.Query(IntervalContainerTest{Id: 1, Low: []int64{1, 1}, High: []int64{1, 1}})
	assert.Len(t, interval_test_query1, 0, "Interval not found: %v", interval_test_query)
}