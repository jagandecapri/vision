package process

import (
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/Workiva/go-datastructures/augmentedtree"
)

type Subspace struct{
	Interval_tree *augmentedtree.Tree
	Grid *tree.Grid
	Subspace_key [2]string
	Scale_factor int
}

func (s *Subspace) ComputeSubspace(mat_old []tree.Point, mat_new_update []tree.Point) {
	subspace_key := s.Subspace_key

	for _, p := range mat_old{
		//fmt.Println("Remove point called")
		rg := s.Grid.GetPointRange(p.Id)
		pnt_container_rem := tree.PointContainer{
			Point: p,
		}
		s.Grid.RemovePoint(pnt_container_rem, rg)
	}
	for _, p := range mat_new_update{
		tmp := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		tmp1 := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		int_container := tree.IntervalContainer{Id: 1, Range: tree.Range{Low: tmp, High: tmp1}, Scale_factor: s.Scale_factor}
		interval := (*s.Interval_tree).Query(int_container)
		if len(interval) > 0{
			interval_ext := interval[0].(tree.IntervalContainer)
			Vec := []float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
			pnt_container := tree.PointContainer{
				Unit_id:  int(interval[0].ID()),
				Vec: Vec,
				Point: p,
			}
			cur_rg := s.Grid.GetPointRange(pnt_container.GetID())
			new_rg := interval_ext.Range
			if cur_rg == (tree.Range{}){
				//fmt.Println("Add point called")
				s.Grid.AddPoint(pnt_container, new_rg)
			} else if cur_rg != (tree.Range{}) && cur_rg != new_rg{
				//fmt.Println("Update point called")
				s.Grid.UpdatePoint(pnt_container, new_rg)
			}
			//fmt.Printf("Interval found %+v %+v \n", int_container, interval)
		} else {
			//fmt.Println("Empty interval:", int_container, interval)
		}
	}
}

func (s *Subspace) Cluster(min_dense_points int, min_cluster_points int){
	s.Grid.Cluster(min_dense_points, min_cluster_points)
}