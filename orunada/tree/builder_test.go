package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRangeBuilder(t *testing.T) {
	min_interval := 0.0
	max_interval := 0.4
	interval_length := 0.1
	ranges := RangeBuilder(min_interval, max_interval, interval_length)
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{0.1, 0.1}}
	r2 := Range{Low: [2]float64{0, 0.1}, High: [2]float64{0.1, 0.2}}
	r3 := Range{Low: [2]float64{0.1, 0}, High: [2]float64{0.2, 0.1}}
	r4 := Range{Low: [2]float64{0.1, 0.1}, High: [2]float64{0.2, 0.2}}
	assert.Contains(t, ranges, r1)
	assert.Contains(t, ranges, r2)
	assert.Contains(t, ranges, r3)
	assert.Contains(t, ranges, r4)
}

func TestIntervalBuilder(t *testing.T) {
	scale_factor := 5
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{0.1, 0.1}}
	r2 := Range{Low: [2]float64{0, 0.1}, High: [2]float64{0.1, 0.2}}
	r3 := Range{Low: [2]float64{0.1, 0}, High: [2]float64{0.2, 0.1}}
	r4 := Range{Low: [2]float64{0.1, 0.1}, High: [2]float64{0.2, 0.2}}
	ranges := []Range{r1,r2,r3,r4}

	i1 := IntervalContainer{Id:0, Range: r1, Scale_factor: scale_factor}
	i2 := IntervalContainer{Id:1, Range: r2, Scale_factor: scale_factor}
	i3 := IntervalContainer{Id:2, Range: r3, Scale_factor: scale_factor}
	i4 := IntervalContainer{Id:3, Range: r4, Scale_factor: scale_factor}
	intervals := IntervalBuilder(ranges, scale_factor)
	assert.Contains(t, intervals, i1)
	assert.Contains(t, intervals, i2)
	assert.Contains(t, intervals, i3)
	assert.Contains(t, intervals, i4)
}