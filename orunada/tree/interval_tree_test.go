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
	}{{interval: IntervalConc{Low: []float64{0.0,0.0}, High: []float64{0.999,0.999}, Decimal_places: 5},
		interval_test: IntervalConc{Low: []float64{0,0.2}, High: []float64{0.998,0.997}, Decimal_places: 5},
		expected: true,
	},
	{interval: IntervalConc{Low: []float64{0.0,0.0}, High: []float64{0.999,0.999}, Decimal_places: 5},
		interval_test: IntervalConc{Low: []float64{1.0,5.0}, High: []float64{1.2,5.5}, Decimal_places: 5},
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
		for _, interval := range intervals{
			interval := augmentedtree.Interval(interval)
			tree.Add(interval)
		}
	})
}