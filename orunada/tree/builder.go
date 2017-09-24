package tree

import (
	"github.com/golang-collections/go-datastructures/augmentedtree"
	"strconv"
	"math/big"
)

func IntervalBuilder(min float64, max float64, interval_length float64, scale_factor int) []IntervalContainer {
	id := 1
	intervals := []IntervalContainer{}
	i := min
	for i < max{
		j := min
		for j < max{
			intervals = append(intervals, IntervalContainer{	Id: id,
				Range: &Range{Low: [2]float64{i, j},
				High: [2]float64{i + interval_length,
					j + interval_length}},
				Scale_factor: scale_factor})
			id += 1
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
	return intervals
}

func NewIntervalTree(dim uint64) augmentedtree.Tree{
	return augmentedtree.New(dim)
}

func NewKDTree(p ...PointInterface) *KDTree{
	return &KDTree{}
}