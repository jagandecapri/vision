package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Create(t *testing.T) {
	assert.Implements(t, (*PointInterface)(nil), new(Unit))
}
