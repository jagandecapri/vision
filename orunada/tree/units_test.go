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
