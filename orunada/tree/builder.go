package tree

import (
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"strconv"
	"math/big"
)

func RangeBuilder(min float64, max float64, interval_length float64) []Range{
	ranges := []Range{}
	i := min
	for i < max{
		j := min
		for j < max{
			tmp := Range{Low: [2]float64{i, j},
				High: [2]float64{i + interval_length,
					j + interval_length}}
			ranges = append(ranges, tmp)
			j += interval_length
			j, _ = strconv.ParseFloat(new(big.Float).SetFloat64(j).Text('f', 1), 64)
			if j >= max{
				break
			}
		}
		i += interval_length
		i, _ = strconv.ParseFloat(new(big.Float).SetFloat64(i).Text('f', 1), 64)
		if i >= max{
			break
		}
	}
	return ranges
}

func IntervalBuilder(ranges []Range,  scale_factor int) []IntervalContainer {
	intervals := []IntervalContainer{}
	for idx, rg := range ranges{
		intervals = append(intervals, IntervalContainer{Id: idx,
			Range: rg,
			Scale_factor: scale_factor})
	}
	return intervals
}

func UnitsBuilder(ranges []Range, dim int) []Unit{
	units := []Unit{}
	for idx, rg := range ranges{
		units = append(units, NewUnit(idx, dim, rg))
	}
	return units
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}

func NewKDTree(p ...PointInterface) *KDTree{
	return &KDTree{}
}