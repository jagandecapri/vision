package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestBuildingKDTree(t *testing.T){
	kd := new(KDTree)


	tests := []struct {
		points []Point
		point_cont *PointContainer
		expected bool
	}{
		{points: []Point{&PointContainer{dim:2, point: []int{3, 6}}, &PointContainer{dim:2, point: []int{17, 15}}, &PointContainer{dim:2, point: []int{13, 15}},
			&PointContainer{dim:2, point: []int{6, 12}}, &PointContainer{dim:2, point: []int{9, 1}}, &PointContainer{dim:2, point: []int{2, 7}}, &PointContainer{dim:2, point: []int{10, 19}}},
			point_cont:  &PointContainer{dim:2, point: []int{3, 6}}, expected: true},
		{points: []Point{&PointContainer{dim:2, point: []int{3, 6}}, &PointContainer{dim:2, point: []int{17, 15}}, &PointContainer{dim:2, point: []int{13, 15}},
			&PointContainer{dim:2, point: []int{6, 12}}, &PointContainer{dim:2, point: []int{9, 1}}, &PointContainer{dim:2, point: []int{2, 7}}, &PointContainer{dim:2, point: []int{10, 19}}},
			point_cont:  &PointContainer{dim:2, point: []int{12, 19}}, expected: false},
	}

	for _, v := range tests{
		kd.Insert(v.points...)
		assert.Equal(t, v.expected, kd.Search(v.point_cont))
	}

	fmt.Println()
}

func TestBFSTreeTraversal(t *testing.T){
	tests := []struct {
		points []Point
		expected [][]int
	}{
		{points: []Point{&PointContainer{dim:2, point: []int{3, 6}}, &PointContainer{dim:2, point: []int{17, 15}}, &PointContainer{dim:2, point: []int{13, 15}},
			&PointContainer{dim:2, point: []int{6, 12}}, &PointContainer{dim:2, point: []int{9, 1}}, &PointContainer{dim:2, point: []int{2, 7}}, &PointContainer{dim:2, point: []int{10, 19}}},
			expected: [][]int{{3,6},{2,7},{17,15},{6,12},{13,15},{9,1},{10,19}}},
	}

	for _, v := range tests{
		kd := new(KDTree)

		kd.Insert(v.points...)
		assert.Equal(t, kd.len, len(v.points))
		res := kd.BFSTraverse()
		for i := 0; i < len(v.expected); i++{
			tmp := res[i].(*PointContainer)
			assert.Equal(t, v.expected[i], tmp.point)
		}
	}
	fmt.Println()
}

func TestKDTree_BFSTraverseChan(t *testing.T) {
	tests := []struct {
		points []Point
		expected [][]int
	}{
		{points: []Point{&PointContainer{dim:2, point: []int{3, 6}}, &PointContainer{dim:2, point: []int{17, 15}}, &PointContainer{dim:2, point: []int{13, 15}},
			&PointContainer{dim:2, point: []int{6, 12}}, &PointContainer{dim:2, point: []int{9, 1}}, &PointContainer{dim:2, point: []int{2, 7}}, &PointContainer{dim:2, point: []int{10, 19}}},
			expected: [][]int{{3,6},{2,7},{17,15},{6,12},{13,15},{9,1},{10,19}}},
	}

	for _,v := range tests{
		kd := new(KDTree)

		kd.Insert(v.points...)
		assert.Equal(t, kd.len, len(v.points))
		out := make(chan Point)
		kd.BFSTraverseChan(out)
		i := 0
		for point := range out{
			tmp := point.(*PointContainer)
			assert.Equal(t, v.expected[i], tmp.point)
			i++
		}
	}
}
