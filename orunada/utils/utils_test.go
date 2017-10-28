package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUniqInt(t *testing.T) {
	input := []int{1,1,2,3,3}
	output := UniqInt(input)
	assert.Contains(t, output, 1)
	assert.Contains(t, output, 2)
	assert.Contains(t, output, 3)
	assert.Equal(t, 3, len(output))
}

func TestUniqInt2(t *testing.T) {
	input := []int{}
	var output []int
	assert.NotPanics(t, func() {
		output = UniqInt(input)
	})
	assert.Equal(t, 0, len(output))
}
