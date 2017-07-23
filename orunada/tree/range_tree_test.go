package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewRangeTree(t *testing.T) {
	dim := uint64(2)
	tr := NewRangeTree(dim)
	assert.NotNil(t, tr)
}
