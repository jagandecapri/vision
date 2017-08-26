package tree

import "github.com/golang-collections/go-datastructures/augmentedtree"

func IntervalBuilder(min float64, max float64, interval_length float64) []IntervalContainer {
	id := 1
	intervals := []IntervalContainer{}
	for i := min; i < max; i += interval_length{
		for j := min; j < max; j += interval_length{
			intervals = append(intervals, IntervalContainer{	Id: id,
				Low: []float64{float64(i), float64(j)},
				High: []float64{float64(i + interval_length), float64(j + interval_length)}})
			id += 1
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