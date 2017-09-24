package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang-collections/go-datastructures/augmentedtree"
)

func TestInterval_OverlapsAtDimension(t *testing.T) {
	tests := []struct {
		interval      IntervalContainer
		interval_test IntervalContainer
		expected      bool
	}{{interval: IntervalContainer{Range: &Range{Low: [2]float64{0.0,0.0}, High: [2]float64{0.999,0.999}}, Scale_factor: 5},
		interval_test: IntervalContainer{Range: &Range{Low: [2]float64{0,0.2}, High: [2]float64{0.998,0.997}}, Scale_factor: 5},
		expected: true,
	},
	{interval: IntervalContainer{Range: &Range{Low: [2]float64{0.0000,0.0000}, High: [2]float64{0.0000,0.0000}}, Scale_factor: 5},
		interval_test: IntervalContainer{Range: &Range{Low: [2]float64{0.0,0.0}, High: [2]float64{1.0,1.0}}, Scale_factor: 5},
		expected: true,
	},
	{interval: IntervalContainer{Range: &Range{Low: [2]float64{0.0,0.0}, High: [2]float64{0.999,0.999}}, Scale_factor: 5},
		interval_test: IntervalContainer{Range: &Range{Low: [2]float64{1.0,5.0}, High: [2]float64{1.2,5.5}}, Scale_factor: 5},
		expected: false,
	}}

	for _, v := range tests{
		assert.Equal(t, v.expected, v.interval.OverlapsAtDimension(v.interval_test, 1), "%v\n", v)
	}
}

func TestInterval_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*augmentedtree.Interval)(nil), new(IntervalContainer))
}

func TestCreateIntervalTree(t *testing.T){
	assert.NotPanics(t, func(){
		min_interval := 0.0
		max_interval := 1.0
		interval_length := 0.1
		scale_factor := 5
		intervals := IntervalBuilder(min_interval, max_interval, interval_length, scale_factor)
		tree := NewIntervalTree(2)
		for _, interval := range intervals{
			interval := augmentedtree.Interval(interval)
			tree.Add(interval)
		}
	})
}

//func TestIntervalQuery(t *testing.T){
//	min_interval := 0.0
//	max_interval := 0.1
//	interval_length := 0.1
//	scale_factor := 5
//	intervals := IntervalBuilder(min_interval, max_interval, interval_length, scale_factor)
//	tree := NewIntervalTree(2)
//	fmt.Println(intervals)
//	for _, interval := range intervals{
//		interval := augmentedtree.Interval(interval)
//		tree.Add(interval)
//	}
//	interval := tree.Query(IntervalContainer{Id: 1, Range{Low: [2]float64{0.0, 0.0}, High: [2]float64{0.0, 0.0}}, Scale_factor: 5})
//	assert.Len(t, interval, 1, "Interval not found: %v", interval)
//	fmt.Println(interval)
//}
//
//
//type IntervalContainerTest struct {
//	Id           int
//	Low          [2]int64
//	High         [2]int64
//}
//
//func (itv IntervalContainerTest) LowAtDimension(dim uint64) int64{
//	tmp :=  itv.Low[dim - 1]
//	return tmp
//}
//
//// HighAtDimension returns an integer representing the higher bound
//// at the requested dimension.
//func (itv IntervalContainerTest) HighAtDimension(dim uint64) int64{
//	tmp := itv.High[dim - 1]
//	return tmp
//}
//
//// OverlapsAtDimension should return a bool indicating if the provided
//// interval overlaps this interval at the dimension requested.
//func (itv IntervalContainerTest) OverlapsAtDimension(interval augmentedtree.Interval, dim uint64) bool{
//	check := false
//	if interval.LowAtDimension(dim) <= itv.HighAtDimension(dim) &&
//		interval.HighAtDimension(dim) >= itv.LowAtDimension(dim){
//		check = true
//	} else {
//		check = false
//	}
//	return check
//}
//
//// ID should be a unique ID representing this interval.  This
//// is used to identify which interval to delete from the tree if
//// there are duplicates.
//func (itv IntervalContainerTest) ID() uint64{
//	return uint64(itv.Id)
//}
//
//func TestIntervalQueryInt(t *testing.T){
//	tree := NewIntervalTree(2)
//	intervals_container := IntervalContainerTest{Id: 1, Low: [2]int64{0, 0}, High: [2]int64{1, 1}}
//	intervals := augmentedtree.Interval(intervals_container)
//	tree.Add(intervals)
//
//	interval_test_query := tree.Query(IntervalContainerTest{Id: 1, Low: [2]int64{0, 0}, High: [2]int64{0, 0}})
//	assert.Len(t, interval_test_query, 1, "Interval not found: %v", interval_test_query)
//	fmt.Println(interval_test_query)
//}