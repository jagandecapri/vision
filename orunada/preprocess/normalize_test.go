package preprocess

import (
	"testing"
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	points := []tree.Point{{Id: 1, Vec_map: map[string]float64{
		"first": 5,
		"second": 10,
	}},
	{Id: 2, Vec_map: map[string]float64{
		"first": 2,
		"second": 6,
	}},
	{Id: 3, Vec_map: map[string]float64{
		"first": 3,
		"second": 7,
	}},
	}

	normalized_points := []tree.Point{{Id: 1, Vec_map: map[string]float64{
		"first": 0.9999999,
		"second": 0.9999999,
	}},
	{Id: 2, Vec_map: map[string]float64{
		"first": 0,
		"second": 0,
	}},
	{Id: 3, Vec_map: map[string]float64{
		"first": 0.3333333333333333,
		"second": 0.25,
	}},
	}

	sorter := []string{"first", "second"}

	norm_points := Normalize(points, sorter)
	assert.Equal(t, normalized_points, norm_points)
}
