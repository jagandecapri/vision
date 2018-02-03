package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCluster_GetCenter(t *testing.T) {

	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	u1 := Unit{Id: 1, Center: Point{Vec: []float64{(0.1 + 0.2 + 0.3 + 0.4 + 0.5)/5, (0.1 + 0.2 + 0.3 + 0.4 + 0.5)/5}},
	Points: map[int]Point{1: {Vec: []float64{0.1, 0.1}},
	2: {Vec: []float64{0.2, 0.2}},
	3: {Vec: []float64{0.3, 0.3}},
	4: {Vec: []float64{0.4, 0.4}},
	5: {Vec: []float64{0.5, 0.5}}},
	Range: r1}

	r2 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	u2 := Unit{Id: 2, Center: Point{Vec: []float64{(1.1 + 1.2 + 1.3 + 1.4)/4, (1.1 + 1.2 + 1.3 + 1.4)/4}},
	Points: map[int]Point{1: {Vec: []float64{1.1, 1.1}},
	2: {Vec: []float64{1.2, 1.2}},
	3: {Vec: []float64{1.3, 1.3}},
	4: {Vec: []float64{1.4, 1.4}}},
	Range: r2}

	cluster := Cluster{Num_of_points: 9, ListOfUnits: map[Range]*Unit{r1: &u1, r2: &u2}}

	x1 := (0.1 + 0.2 + 0.3 + 0.4 + 0.5 + 1.1 + 1.2 + 1.3 + 1.4)/9
	x2 := (0.1 + 0.2 + 0.3 + 0.4 + 0.5 + 1.1 + 1.2 + 1.3 + 1.4)/9
	center := Point{Vec: []float64{x1, x2}}

	pc := cluster.GetCenter()

	assert.Equal(t, center, pc)
}
