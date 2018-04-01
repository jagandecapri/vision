package tree

import (
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/jagandecapri/vision/utils"
)

func RangeBuilder(min float64, max float64, interval_length float64) []Range{
	ranges := []Range{}
	i := min
	for i < max{
		j := min
		for j < max{
			tmp_i := utils.Round(i + interval_length, 0.1)
			tmp_j := utils.Round(j + interval_length, 0.1)
			tmp := Range{Low: [2]float64{i, j},
				High: [2]float64{tmp_i,
					tmp_j}}
			ranges = append(ranges, tmp)
			j = utils.Round(j + interval_length, 0.1)
			if j >= max{
				break
			}
		}
		i = utils.Round(i + interval_length, 0.1)
		if i >= max{
			break
		}
	}
	return ranges
}

func IntervalBuilder(ranges []Range,  scale_factor int) map[Range]IntervalContainer {
	intervals := make(map[Range]IntervalContainer)
	for idx, rg := range ranges{
		intervals[rg] = IntervalContainer{Id: idx,
			Range: rg,
			Scale_factor: scale_factor}
	}
	return intervals
}

func UnitsBuilder(ranges []Range, dim int) map[Range]*Unit{
	units := make(map[Range]*Unit)
	for idx, rg := range ranges{
		units[rg] = NewUnit(idx, dim, rg)
	}
	return units
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}