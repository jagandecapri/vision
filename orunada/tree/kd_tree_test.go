package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestBuildingKDTree(t *testing.T){
	kd := new(KDTree)

	points := [][]int{{3, 6}, {17, 15}, {13, 15}, {6, 12}, {9, 1}, {2, 7}, {10, 19}}

	for i := 0; i < len(points); i++ {
		kd.Insert(points[i]...)
	}

	tests := []struct {
		point []int
		expected bool
	}{
		{point:  []int{10, 19}, expected: true},
		{point:  []int{12, 19}, expected: false},
	}

	for _, v := range tests{
		assert.Equal(t, v.expected, kd.Search(v.point...))
	}
}

func TestBFSTreeTraversal(t *testing.T){
	tests := []struct {
		points [][]int
		expected [][]int
	}{
		{points: [][]int{{3, 6}, {17, 15}, {13, 15}, {6, 12}, {9, 1}, {2, 7}, {10, 19}},
			expected: [][]int{{3,6},{2,7},{17,15},{6,12},{13,15},{9,1},{10,19}}},
	}

	for _, v := range tests{
		kd := new(KDTree)

		for i := 0; i < len(v.points); i++ {
			kd.Insert(v.points[i]...)
		}
		assert.Equal(t, kd.len, len(v.points))
		res := kd.BFSTraverse()
		for i := 0; i < len(v.expected); i++{
			assert.Equal(t, v.expected[i], res[i])
		}
	}
	fmt.Println()
}
