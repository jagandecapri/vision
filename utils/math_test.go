package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

func TestGaussian(t *testing.T) {
	res := Gaussian(1,1,1,1)
	assert.Equal(t, 1.0, res)

	res = Gaussian(1,1,1,0)
	assert.True(t, math.IsNaN(res))

	res = Gaussian(2, 1, 0, 1)
	assert.Equal(t, 1.4142135623731, res)
}
