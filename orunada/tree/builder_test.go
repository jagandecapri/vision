package tree

import (
	"testing"
)

func TestIntervalBuilder(t *testing.T) {
	min_interval := 0.0
	max_interval := 1.0
	interval_length := 0.1
	scale_factor := 5
	IntervalBuilder(min_interval, max_interval, interval_length, scale_factor)
	//Add assertion
}
