package tree

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T){
	points := []point{{x: 0.1, y: 0.1, z: 0.1}, {x: 0.2, y: 0.2, z: 0.2}, {x: 0.3, y: 0.3, z: 0.3}, {x: 0.4, y: 0.4, z: 0.4}}
	Main_mock(points, 0.3, 2)
	assert.False(t, true)
}