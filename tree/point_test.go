package tree

import (
	"testing"
	"math"
	"github.com/stretchr/testify/assert"
)

func TestPointContainer_Distance(t *testing.T) {
	tests := []struct{
		point1 *Point
		point2 *Point
		expected float64
	}{
		{point1: &Point{Id: 1,
			Unit_id: 1,
			Vec: []float64{1.0,2.0},
			Vec_map: map[string]float64{"a": 1.0, "b": 2.0},
			},
		point2: &Point{Id: 2,
			Unit_id: 2,
			Vec: []float64{3.0,4.0},
			Vec_map: map[string]float64{"a": 3.0, "b": 4.0},
		},
			expected: math.Sqrt(8),
		},
	}

	for _, v := range tests{
		res := v.point1.Distance(v.point2)
		assert.Equal(t, v.expected, res)
	}
}
