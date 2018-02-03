package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGrid_AddUnit(t *testing.T) {
	tmp := Range{Low: [2]float64{0.0, 0.0}, High: [2]float64{1.0, 1.0}}
	unit := &Unit{Range: tmp}
	grid := NewGrid()
	grid.AddUnit(unit)
}

func TestGrid_GetNeighbouringUnits(t *testing.T) {
	grid := Grid{Store: map[Range]*Unit{{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1},
		{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2},
		{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}: {Id: 3},
		{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}: {Id: 4},
		{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}: {Id: 5},
		{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}: {Id: 6},
		{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}: {Id: 7},
		{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}: {Id: 8},
		{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}: {Id: 9}},
	}

	neighbouring_units := grid.GetNeighbouringUnits(Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}, 1.0)
	unit_ids := []int{}
	for _, nu := range neighbouring_units{
		unit_ids = append(unit_ids, nu.Id)
	}

	for i := 1; i <= 8; i++{
		assert.Contains(t, unit_ids, i)
	}

}

func TestGrid_AddPoint(t *testing.T) {
	p := Point{Id: 1}
	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	grid := NewGrid()
	unit := NewUnit(1, 2, rg)
	grid.Store[rg] = &unit
	grid.AddPoint(p, rg)
	assert.Equal(t, grid.point_unit_map[1], rg)
	assert.Equal(t, grid.Store[rg].Points[1].Id, 1)
	assert.False(t, grid.Store[rg].Center_calculated)

	p1 := Point{Id: 2}
	rg1 := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	assert.NotPanics(t, func(){grid.AddPoint(p1, rg1)})
}

func TestGrid_RemovePoint(t *testing.T) {
	p := Point{Id: 1}
	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	grid := NewGrid()
	unit := NewUnit(1, 2, rg)
	grid.Store[rg] = &unit
	grid.AddPoint(p, rg)
	grid.RemovePoint(p, rg)
	var ok bool
	_, ok = grid.point_unit_map[1]
	assert.False(t, ok)
	_, ok = grid.Store[rg].Points[1]
	assert.False(t, ok)
	assert.False(t, grid.Store[rg].Center_calculated)

	p1 := Point{Id: 2}
	rg1 := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	assert.NotPanics(t, func(){grid.RemovePoint(p1, rg1)})
}

func TestGrid_UpdatePoint(t *testing.T) {
	p := Point{Id: 1}
	cur_rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	grid := NewGrid()
	unit := NewUnit(1, 2, cur_rg)
	grid.Store[cur_rg] = &unit
	grid.AddPoint(p, cur_rg)
	new_rg := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	unit1 := NewUnit(1, 2, new_rg)
	grid.Store[new_rg] = &unit1
	grid.UpdatePoint(p, new_rg)
	var ok bool
	_, ok = grid.point_unit_map[1]
	assert.True(t, ok)
	_, ok = grid.Store[cur_rg].Points[1]
	assert.False(t, ok)
	_, ok = grid.Store[new_rg].Points[1]
	assert.True(t, ok)
}

func TestGrid_RecomputeDenseUnits(t *testing.T) {
	p := Point{Id: 1}

	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	unit := NewUnit(1,2,rg)
	unit.Points[1] = p

	rg1 := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	unit1 := NewUnit(1,2,rg1)

	grid := NewGrid()
	grid.Store[rg] = &unit
	grid.Store[rg1] = &unit1

	var ok bool
	var listNewDenseUnits, listOldDenseUnits map[Range]*Unit
	/*
	if isDenseUnit(unit, min_dense_points){
		.....
		if !ok{
			us.listDenseUnits[rg] = unit
			us.listNewDenseUnits[rg] = unit
	}
	 */
	listNewDenseUnits, _ = grid.RecomputeDenseUnits(1)
	_, ok =  grid.listDenseUnits[rg]
	assert.True(t, ok)
	_, ok = listNewDenseUnits[rg]
	assert.True(t, ok)
	_, ok = grid.listDenseUnits[rg1]
	assert.False(t, ok)

	/*
	if isDenseUnit(unit, min_dense_points){
		_, ok := us.listDenseUnits[rg]
		...
	}
	 */
	listNewDenseUnits, _ = grid.RecomputeDenseUnits(1)
	_, ok =  grid.listDenseUnits[rg]
	assert.True(t, ok)

	/*
	 else {
		...
		if ok{
			delete(us.listDenseUnits, rg)
			us.listOldDenseUnits[rg] = unit
		}
	 */
	unit.Points = make(map[int]Point)
	_, listOldDenseUnits = grid.RecomputeDenseUnits(1)
	_, ok =  grid.listDenseUnits[rg]
	assert.False(t, ok)
	_, ok = listOldDenseUnits[rg]
	assert.True(t, ok)

	/*
	else {
		_, ok := us.listDenseUnits[rg]

	 */
	_, listOldDenseUnits =grid.RecomputeDenseUnits(1)
	_, ok = listOldDenseUnits[rg1]
	assert.False(t, ok)
}

func TestGrid_ProcessOldDenseUnits(t *testing.T) {
	rg := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	unit := NewUnit(1,2,rg)
	unit.Id = 9
	unit.Cluster_id = 0
	unit.neighbour_units = map[Range]*Unit{{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1, Cluster_id: 1},
		{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2, Cluster_id: 1},
		{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}: {Id: 3},
		{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}: {Id: 4},
		{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}: {Id: 5},
		{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}: {Id: 6},
		{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}: {Id: 7},
		{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}: {Id: 8}}
	grid := NewGrid()
	list_of_units := map[Range]*Unit{
		Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1, Cluster_id: 1},
		Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2, Cluster_id: 1},
	}

	c := Cluster{Cluster_id:0, ListOfUnits: list_of_units}
	grid.AddUpdateCluster(c)

	listOldDenseUnits := map[Range]*Unit{{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}: &unit}

	listUnitToRep := grid.ProcessOldDenseUnits(listOldDenseUnits)

	unit_ids := []int{}
	for _, nu := range listUnitToRep{
		unit_ids = append(unit_ids, nu.Id)
	}

	assert.Contains(t, unit_ids, 1)
	assert.Contains(t, unit_ids, 2)
}

func TestGrid_GetPointRange(t *testing.T) {
	grid := NewGrid()
	rg := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	grid.point_unit_map[1] = rg

	assert.Equal(t, grid.GetPointRange(1), rg)
	assert.Equal(t, grid.GetPointRange(2), Range{})
}