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
		{point1: &PointContainer{dim: 2, point: []int{1,2}},
			point2: &PointContainer{dim: 2, point: []int{3,4}},
			expected: math.Sqrt(8),
		},
	}

	for _, v := range tests{
		res := v.point1.Distance(v.point2)
		assert.Equal(t, v.expected, res)
	}
}
