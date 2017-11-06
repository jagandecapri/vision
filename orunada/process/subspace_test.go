package process

import (
	"testing"
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/stretchr/testify/assert"
)

func TestSubspace_ComputeSubspace(t *testing.T) {
	subspace_key := [2]string{"first", "second"}
	min_interval := 0.0
	max_interval := 0.2
	interval_length := 0.1
	dim := 2
	scale_factor := 5
	ranges := tree.RangeBuilder(min_interval, max_interval, interval_length)
	intervals := tree.IntervalBuilder(ranges, scale_factor)
	units := tree.UnitsBuilder(ranges, dim)

	int_tree := tree.NewIntervalTree(uint64(dim))
	Unit := tree.NewGrid()

	for _, interval := range intervals{
		int_tree.Add(interval)
	}

	for rg, unit := range units{
		Unit.AddUnit(&unit, rg)
	}

	subspace := Subspace{Interval_tree: &int_tree,
		Grid: &Unit,
		Subspace_key: subspace_key,
		Scale_factor: scale_factor,
	}

	//To test new points
	p1 := tree.Point{Id: 1, Vec_map: map[string]float64{
		"first": 0.05,
		"second": 0.05,
		"three": 0.15,
	}}
	p2 := tree.Point{Id: 2, Vec_map: map[string]float64{
		"first": 0.15,
		"second": 0.15,
		"three": 0.15,
	}}

	points := []tree.Point{p1, p2}
	subspace.ComputeSubspace([]tree.Point{}, points)

	assert.Equal(t, tree.Range{Low: [2]float64{0, 0}, High: [2]float64{0.1, 0.1}}, subspace.Grid.Point_unit_map[1])
	assert.Equal(t, tree.Range{Low: [2]float64{0.1, 0.1}, High: [2]float64{0.2, 0.2}}, subspace.Grid.Point_unit_map[2])

	//To test updating points
	p1 = tree.Point{Id: 1, Vec_map: map[string]float64{
		"first": 0.15,
		"second": 0.15,
		"three": 0.15,
	}}
	p2 = tree.Point{Id: 2, Vec_map: map[string]float64{
		"first": 0.05,
		"second": 0.05,
		"three": 0.15,
	}}

	points = []tree.Point{p1, p2}
	subspace.ComputeSubspace([]tree.Point{}, points)

	assert.Equal(t, tree.Range{Low: [2]float64{0.1, 0.1}, High: [2]float64{0.2, 0.2}}, subspace.Grid.Point_unit_map[1])
	assert.Equal(t, tree.Range{Low: [2]float64{0, 0}, High: [2]float64{0.1, 0.1}}, subspace.Grid.Point_unit_map[2])

	//To test removing points
	p1 = tree.Point{Id: 1, Vec_map: map[string]float64{
		"first": 0.15,
		"second": 0.15,
		"three": 0.15,
	}}
	p2 = tree.Point{Id: 2, Vec_map: map[string]float64{
		"first": 0.05,
		"second": 0.05,
		"three": 0.15,
	}}

	points = []tree.Point{p1, p2}
	subspace.ComputeSubspace(points, []tree.Point{})

	assert.Equal(t, tree.Range{}, subspace.Grid.Point_unit_map[1])
	assert.Equal(t, tree.Range{}, subspace.Grid.Point_unit_map[2])
}