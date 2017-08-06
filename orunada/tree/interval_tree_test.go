package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang-collections/go-datastructures/augmentedtree"
)

func TestInterval_OverlapsAtDimension(t *testing.T) {
	tests := []struct {
		interval      IntervalConc
		interval_test IntervalConc
		expected      bool
	}{{interval: IntervalConc{low: []int64{0,0}, high: []int64{999,999}},
		interval_test: IntervalConc{low: []int64{0,2}, high: []int64{998,997}},
		expected: true,
	},
	{interval: IntervalConc{low: []int64{0,0}, high: []int64{999,999}},
		interval_test: IntervalConc{low: []int64{0,2}, high: []int64{999,999}},
		expected: false,
	}}

	for _, v := range tests{
		assert.Equal(t, v.expected, v.interval.OverlapsAtDimension(v.interval_test, 1), "%v\n", v)
	}
}

func TestInterval_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*augmentedtree.Interval)(nil), new(IntervalConc))
}

func TestCreateIntervalTree(t *testing.T){
	assert.NotPanics(t, func(){
		intervals := IntervalBuilder(0, 10, 1)
		tree := NewIntervalTree(2)
		tree.Add(intervals...)
	})
}