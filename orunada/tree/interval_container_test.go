package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"fmt"
)

func TestInterval_OverlapsAtDimension(t *testing.T) {
	tests := []struct {
		interval      IntervalContainer
		interval_test IntervalContainer
		expected      bool
	}{{interval: IntervalContainer{Low: []float64{0.0,0.0}, High: []float64{0.999,0.999}, Scale_factor: 5},
		interval_test: IntervalContainer{Low: []float64{0,0.2}, High: []float64{0.998,0.997}, Scale_factor: 5},
		expected: true,
	},
	{interval: IntervalContainer{Low: []float64{0.0000,0.0000}, High: []float64{0.0000,0.0000}, Scale_factor: 5},
		interval_test: IntervalContainer{Low: []float64{0.0,0.0}, High: []float64{1.0,1.0}, Scale_factor: 5},
		expected: true,
	},
	{interval: IntervalContainer{Low: []float64{0.0,0.0}, High: []float64{0.999,0.999}, Scale_factor: 5},
		interval_test: IntervalContainer{Low: []float64{1.0,5.0}, High: []float64{1.2,5.5}, Scale_factor: 5},
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
		intervals := IntervalBuilder(0.0, 1.0, 0.1)
		tree := NewIntervalTree(2)
		fmt.Println(intervals)
		for _, interval := range intervals{
			interval := augmentedtree.Interval(interval)
			tree.Add(interval)
		}
		interval := tree.Query(IntervalContainer{Id: 1, Low: []float64{0, 0}, High: []float64{0, 0}})
		fmt.Println(interval)
	})
}