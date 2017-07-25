package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInterval_OverlapsAtDimension(t *testing.T) {
	tests := []struct {
		interval Interval
		interval_test Interval
		expected bool
	}{{interval: Interval{low: []int64{0,0}, high: []int64{999,999}},
		interval_test: Interval{low: []int64{0,2}, high: []int64{998,997}},
		expected: true,
	},
	{interval: Interval{low: []int64{0,0}, high: []int64{999,999}},
		interval_test: Interval{low: []int64{0,2}, high: []int64{999,999}},
		expected: false,
	}}

	for _, v := range tests{
		assert.Equal(t, v.expected, v.interval.OverlapsAtDimension(v.interval_test, 1), "%v\n", v)
	}
}
