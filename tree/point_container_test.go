package tree

import (
	"testing"
	"math"
	"github.com/stretchr/testify/assert"
)

func TestPointContainer_Distance(t *testing.T) {
	tests := []struct{
		point1 *PointContainer
		point2 *PointContainer
		expected float64
	}{
		{point1: &PointContainer{1,
			[]float64{1.0,2.0},
			Point{1,
				map[string]float64{"a": 1.0, "b": 2.0},
			}},
		point2: &PointContainer{2,
			[]float64{3.0,4.0},
			Point{2,
				map[string]float64{"a": 3.0, "b": 4.0},
		}},
			expected: math.Sqrt(8),
		},
	}

	for _, v := range tests{
		res := v.point1.Distance(v.point2)
		assert.Equal(t, v.expected, res)
	}
}
