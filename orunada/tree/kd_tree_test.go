package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

type PointContainer struct{
	dim int
	point []int
}

func (p PointContainer) Dim() int{
	return p.dim
}

func (p PointContainer) GetValue(dim int) int{
	return p.point[dim]
}

func (p *PointContainer) Distance(point Point) float64{
	return 0.0
}

func (p *PointContainer) PlaneDistance(val float64, dim int) float64{
	return 0.0
}

func TestBuildingKDTree(t *testing.T){
	kd := new(KDTree)

	points := []Point{&PointContainer{dim:2, point: []int{3, 6}}, &PointContainer{dim:2, point: []int{17, 15}}, &PointContainer{dim:2, point: []int{13, 15}},
		&PointContainer{dim:2, point: []int{6, 12}}, &PointContainer{dim:2, point: []int{9, 1}}, &PointContainer{dim:2, point: []int{2, 7}}, &PointContainer{dim:2, point: []int{10, 19}}}

	for i := 0; i < len(points); i++ {
		kd.Insert(points[i])
	}

	tests := []struct {
		point_cont *PointContainer
		expected bool
	}{
		{point_cont:  &PointContainer{dim:2, point: []int{3, 6}}, expected: true},
		{point_cont:  &PointContainer{dim:2, point: []int{12, 19}}, expected: false},
	}

	for _, v := range tests{
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

		for i := 0; i < len(v.points); i++ {
			kd.Insert(v.points[i])
		}
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

		for i := 0; i < len(v.points); i++ {
			kd.Insert(v.points[i])
		}
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
