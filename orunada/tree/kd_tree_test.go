package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"fmt"
)

func TestKDTree_InsertTest(t *testing.T){
	tests := []struct {
		points []PointInterface
	}{
		{points: []PointInterface{&PointContainer{1,
			[]float64{1.0, 2.0},
			Point{1,
				map[string]float64{"a": 1.0, "b": 2.0},
			},
		},
		}},
	}

	for _, v := range tests{
		kd := NewKDTree()
		kd.Insert(v.points...)
		assert.Equal(t, kd.len, len(v.points))
	}
}

func TestKDTree_Insert(t *testing.T) {
	tests := []struct {
		points []PointInterface
		point_cont *PointContainer
		expected bool
	}{
		{points: []PointInterface{&PointContainer{1,
			[]float64{1.0, 2.0},
			Point{1,
				map[string]float64{"a": 1.0, "b": 2.0},
			},
		}},
			point_cont:  &PointContainer{1,
				[]float64{1.0, 2.0},
				Point{1,
					map[string]float64{"a": 1.0, "b": 2.0},
				},
			},
			expected: true},
		{points: []PointInterface{&PointContainer{1,
			[]float64{1.0, 2.0},
			Point{1,
				map[string]float64{"a": 1.0, "b": 2.0},
			},
		}},
			point_cont:  &PointContainer{1,
				[]float64{3.0, 2.0},
				Point{1,
					map[string]float64{"a": 3.0, "b": 2.0},
				},
			},
			expected: false},
	}

	for _, v := range tests{
		kd := NewKDTree()
		kd.Insert(v.points...)
		assert.Equal(t, v.expected, kd.Search(v.point_cont))
	}

	fmt.Println()
}

func TestKDTree_BFSTraverseTest(t *testing.T){
	tests := []struct {
		points []PointInterface
		expected [][]float64
	}{
		{points: []PointInterface{&PointContainer{1,
			[]float64{3.0, 6.0},
			Point{1,
				map[string]float64{"a": 3.0, "b": 6.0},
			},
		},
		&PointContainer{1,
			[]float64{6.0, 12.0},
			Point{1,
				map[string]float64{"a": 6.0, "b": 12.0},
			},
		},
		&PointContainer{1,
			[]float64{2.0, 7.0},
			Point{1,
				map[string]float64{"a": 2.0, "b": 7.0},
			},
		},
		&PointContainer{1,
			[]float64{13.0, 15.0},
			Point{1,
				map[string]float64{"a": 13.0, "b": 15.0},
			},
		},
		&PointContainer{1,
			[]float64{17.0, 15.0},
			Point{1,
				map[string]float64{"a": 17.0, "b": 15.0},
			},
		},
		&PointContainer{1,
			[]float64{9.0, 1.0},
			Point{1,
				map[string]float64{"a": 9.0, "b": 1.0},
			},
		},
		&PointContainer{1,
			[]float64{10.0, 19.0},
			Point{1,
				map[string]float64{"a": 10.0, "b": 19.0},
			},
		},
		},
			expected: [][]float64{{3.0,6.0},{2.0,7.0},{6.0,12.0},{9.0,1.0},{13.0,15.0},{10.0,19.0},{17.0,15.0}}},
	}

	for _, v := range tests{
		kd := NewKDTree()
		kd.Insert(v.points...)
		assert.Equal(t, kd.len, len(v.points))
		res := kd.BFSTraverse()
		for i := 0; i < len(v.expected); i++{
			tmp := res[i].(*PointContainer)
			assert.Equal(t, v.expected[i], tmp.Vec)
		}
	}
	fmt.Println()
}

func TestKDTree_BFSTraverseChan(t *testing.T) {
	tests := []struct {
		points []PointInterface
		expected [][]float64
	}{
		{points: []PointInterface{&PointContainer{1,
			[]float64{3.0, 6.0},
			Point{1,
				map[string]float64{"a": 3.0, "b": 6.0},
			},
		},
			&PointContainer{1,
				[]float64{6.0, 12.0},
				Point{1,
					map[string]float64{"a": 6.0, "b": 12.0},
				},
			},
			&PointContainer{1,
				[]float64{2.0, 7.0},
				Point{1,
					map[string]float64{"a": 2.0, "b": 7.0},
				},
			},
			&PointContainer{1,
				[]float64{13.0, 15.0},
				Point{1,
					map[string]float64{"a": 13.0, "b": 15.0},
				},
			},
			&PointContainer{1,
				[]float64{17.0, 15.0},
				Point{1,
					map[string]float64{"a": 17.0, "b": 15.0},
				},
			},
			&PointContainer{1,
				[]float64{9.0, 1.0},
				Point{1,
					map[string]float64{"a": 9.0, "b": 1.0},
				},
			},
			&PointContainer{1,
				[]float64{10.0, 19.0},
				Point{1,
					map[string]float64{"a": 10.0, "b": 19.0},
				},
			},
		},
			expected: [][]float64{{3.0,6.0},{2.0,7.0},{6.0,12.0},{9.0,1.0},{13.0,15.0},{10.0,19.0},{17.0,15.0}}},
	}

	for _,v := range tests{
		kd := NewKDTree()
		kd.Insert(v.points...)
		assert.Equal(t, kd.len, len(v.points))
		out := make(chan PointInterface)
		kd.BFSTraverseChan(out)
		i := 0
		for point := range out{
			tmp := point.(*PointContainer)
			assert.Equal(t, v.expected[i], tmp.Vec)
			i++
		}
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
