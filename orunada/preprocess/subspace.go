package preprocess

import (
	"github.com/jagandecapri/vision/orunada/tree"
	"github.com/Workiva/go-datastructures/augmentedtree"
)

type Subspace struct{
	Interval_tree *augmentedtree.Tree
	Units *tree.Units
	Subspace_key [2]string
	Scale_factor int
}

func (s *Subspace) ComputeSubspace(mat []tree.Point) {
	subspace_key := s.Subspace_key
	for _, p := range mat{
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
			cur_rg := s.Units.GetPointRange(pnt_container.GetID())
			new_rg := interval_ext.Range
			if cur_rg == (tree.Range{}){
				s.Units.AddPoint(pnt_container, new_rg)
			} else if cur_rg != (tree.Range{}) && cur_rg != new_rg{
				s.Units.UpdatePoint(pnt_container, new_rg)
			}
		} else {
		}
	}
}