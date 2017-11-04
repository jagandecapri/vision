package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUnits_AddUnit(t *testing.T) {
	tmp := Range{Low: [2]float64{0.0, 0.0}, High: [2]float64{1.0, 1.0}}
	unit := &Unit{}
	units := NewUnits()
	units.AddUnit(unit, tmp)
}

func TestUnits_GetNeighbouringUnits(t *testing.T) {
	units := Units{Store: map[Range]*Unit{{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1},
		{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2},
		{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}: {Id: 3},
		{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}: {Id: 4},
		{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}: {Id: 5},
		{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}: {Id: 6},
		{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}: {Id: 7},
		{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}: {Id: 8},
		{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}: {Id: 9}},
	}

	neighbouring_units := units.GetNeighbouringUnits(Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}, 1.0)
	unit_ids := []int{}
	for _, nu := range neighbouring_units{
		unit_ids = append(unit_ids, nu.Id)
	}

	for i := 1; i <= 8; i++{
		assert.Contains(t, unit_ids, i)
	}

}

func TestUnits_ImplementsClusterInterface(t *testing.T){
	assert.Implements(t, (*ClusterInterface)(nil), new(Units))
}

func TestUnits_AddPoint(t *testing.T) {
	p := PointContainer{Point: Point{Id: 1}}
	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	units := NewUnits()
	unit := NewUnit(1, 2, rg)
	units.Store[rg] = &unit
	units.AddPoint(p, rg)
	assert.Equal(t, units.Point_unit_map[1], rg)
	assert.Equal(t, units.Store[rg].points[1].Id, 1)
	assert.False(t, units.Store[rg].Center_calculated)
}

func TestUnits_RemovePoint(t *testing.T) {
	p := PointContainer{Point: Point{Id: 1}}
	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	units := NewUnits()
	unit := NewUnit(1, 2, rg)
	units.Store[rg] = &unit
	units.AddPoint(p, rg)
	units.RemovePoint(p, rg)
	var ok bool
	_, ok = units.Point_unit_map[1]
	assert.False(t, ok)
	_, ok = units.Store[rg].points[1]
	assert.False(t, ok)
	assert.False(t, units.Store[rg].Center_calculated)
}

func TestUnits_UpdatePoint(t *testing.T) {
	p := PointContainer{Point: Point{Id: 1}}
	cur_rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	units := NewUnits()
	unit := NewUnit(1, 2, cur_rg)
	units.Store[cur_rg] = &unit
	units.AddPoint(p, cur_rg)
	new_rg := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	unit1 := NewUnit(1, 2, new_rg)
	units.Store[new_rg] = &unit1
	units.UpdatePoint(p, new_rg)
	var ok bool
	_, ok = units.Point_unit_map[1]
	assert.True(t, ok)
	_, ok = units.Store[cur_rg].points[1]
	assert.False(t, ok)
	_, ok = units.Store[new_rg].points[1]
	assert.True(t, ok)
}

func TestUnits_RecomputeDenseUnits(t *testing.T) {
	p := PointContainer{Point : Point{Id: 1}}

	rg := Range{Low: [2]float64{0, 0}, High: [2]float64{0.5, 0.5}}
	unit := NewUnit(1,2,rg)
	unit.points[1] = p

	rg1 := Range{Low: [2]float64{0.5, 0.5}, High: [2]float64{1.0, 1.0}}
	unit1 := NewUnit(1,2,rg1)

	units := NewUnits()
	units.Store[rg] = &unit
	units.Store[rg1] = &unit1

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
	listNewDenseUnits, _ = units.RecomputeDenseUnits(1)
	_, ok =  units.listDenseUnits[rg]
	assert.True(t, ok)
	_, ok = listNewDenseUnits[rg]
	assert.True(t, ok)
	_, ok = units.listDenseUnits[rg1]
	assert.False(t, ok)

	/*
	if isDenseUnit(unit, min_dense_points){
		_, ok := us.listDenseUnits[rg]
		...
	}
	 */
	listNewDenseUnits, _ = units.RecomputeDenseUnits(1)
	_, ok =  units.listDenseUnits[rg]
	assert.True(t, ok)

	/*
	 else {
		...
		if ok{
			delete(us.listDenseUnits, rg)
			us.listOldDenseUnits[rg] = unit
		}
	 */
	unit.points = make(map[int]PointContainer)
	_, listOldDenseUnits = units.RecomputeDenseUnits(1)
	_, ok =  units.listDenseUnits[rg]
	assert.False(t, ok)
	_, ok = listOldDenseUnits[rg]
	assert.True(t, ok)

	/*
	else {
		_, ok := us.listDenseUnits[rg]

	 */
	_, listOldDenseUnits =units.RecomputeDenseUnits(1)
	_, ok = listOldDenseUnits[rg1]
	assert.False(t, ok)
}

func TestUnits_ProcessOldDenseUnits(t *testing.T) {
	rg := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	unit := NewUnit(1,2,rg)
	unit.Id = 9
	unit.Cluster_id = 1
	unit.Neighbour_units = map[Range]*Unit{{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1, Cluster_id: 1},
		{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2, Cluster_id: 1},
		{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}: {Id: 3},
		{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}: {Id: 4},
		{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}: {Id: 5},
		{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}: {Id: 6},
		{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}: {Id: 7},
		{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}: {Id: 8}}
	units := NewUnits()
	list_of_units := map[Range]*Unit{
		Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}: {Id: 1, Cluster_id: 1},
		Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}: {Id: 2, Cluster_id: 1},
	}
	units.Cluster_map = map[int]Cluster{1: {ListOfUnits: list_of_units}}
	listOldDenseUnits := map[Range]*Unit{{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}: &unit}

	listUnitToRep := units.ProcessOldDenseUnits(listOldDenseUnits)

	unit_ids := []int{}
	for _, nu := range listUnitToRep{
		unit_ids = append(unit_ids, nu.Id)
	}

	assert.Contains(t, unit_ids, 1)
	assert.Contains(t, unit_ids, 2)
}

func TestUnits_GetPointRange(t *testing.T) {
	units := NewUnits()
	rg := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.Point_unit_map[1] = rg

	assert.Equal(t, units.GetPointRange(1), rg)
	assert.Equal(t, units.GetPointRange(2), Range{})
}