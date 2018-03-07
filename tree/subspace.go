package tree

import (
	"github.com/Workiva/go-datastructures/augmentedtree"
)

type Subspace struct{
	interval_tree *augmentedtree.Tree
	*Grid
	Subspace_key  [2]string
	Scale_factor  int
}

func (s *Subspace) GetIntervalTree()*augmentedtree.Tree{
	return s.interval_tree
}

func (s *Subspace) SetIntervalTree(interval_tree *augmentedtree.Tree){
	s.interval_tree = interval_tree
}

func (s *Subspace) ComputeSubspace(mat_old []Point, mat_new_update []Point) {
	subspace_key := s.Subspace_key

	for _, p := range mat_old{
		//fmt.Println("Remove point called")
		point := Point{Id: p.Id}
		rg := s.Grid.GetPointRange(point.Id)
		s.Grid.RemovePoint(point, rg) //Needs only ID for deletion
	}
	for _, p := range mat_new_update{
		subspace_key_0 := subspace_key[0]
		subspace_key_1 := subspace_key[1]
		subspace_val_0 := p.Vec_map[subspace_key_0]
		subspace_val_1 := p.Vec_map[subspace_key_1]
		Vec := [2]float64{subspace_val_0, subspace_val_1}
		int_container := IntervalContainer{Id: 1, Range: Range{Low: Vec, High: Vec}, Scale_factor: s.Scale_factor}
		interval := (*s.interval_tree).Query(int_container)
		if len(interval) > 0{
			interval_ext := interval[0].(IntervalContainer)
			p.Unit_id = int(interval[0].ID())
			cur_rg := s.Grid.GetPointRange(p.GetID())
			new_rg := interval_ext.Range
			if cur_rg == (Range{}){
				//fmt.Println("Add point called")
				s.Grid.AddPoint(p, new_rg)
			} else if cur_rg != (Range{}) && cur_rg != new_rg{
				//fmt.Println("Update point called")
				s.Grid.UpdatePoint(p, new_rg)
			}
			//fmt.Printf("Interval found %+v %+v \n", int_container, interval)
		} else {
			//fmt.Println("Empty interval:", int_container, interval)
		}
	}
}

//func (s *Subspace) ComputeSubspace(mat_old []Point, mat_new_update []Point) {
//	subspace_key := s.Subspace_key
//
//	for _, p := range mat_old{
//		//fmt.Println("Remove point called")
//		point := Point{Id: p.Id}
//		rg := s.Grid.GetPointRange(point.Id)
//		s.Grid.RemovePoint(point, rg) //Needs only ID for deletion
//	}
//	for _, p := range mat_new_update{
//		subspace_key_0 := subspace_key[0]
//		subspace_key_1 := subspace_key[1]
//		subspace_val_0 := p.Vec_map[subspace_key_0]
//		subspace_val_1 := p.Vec_map[subspace_key_1]
//		Vec := [2]float64{subspace_val_0, subspace_val_1}
//		int_container := IntervalContainer{Id: 1, Range: Range{Low: Vec, High: Vec}, Scale_factor: s.Scale_factor}
//		interval := (*s.interval_tree).Query(int_container)
//		if len(interval) > 0{
//			interval_ext := interval[0].(IntervalContainer)
//			Vec_map := map[string]float64{subspace_key_0: subspace_val_0, subspace_key_1: subspace_val_1}
//			point := Point{
//				Id: p.Id,
//				Key: p.Key,
//				Unit_id:  int(interval[0].ID()),
//				Vec: Vec[:],
//				Vec_map:  Vec_map,
//			}
//
//			cur_rg := s.Grid.GetPointRange(point.GetID())
//			new_rg := interval_ext.Range
//			if cur_rg == (Range{}){
//				//fmt.Println("Add point called")
//				s.Grid.AddPoint(point, new_rg)
//			} else if cur_rg != (Range{}) && cur_rg != new_rg{
//				//fmt.Println("Update point called")
//				s.Grid.UpdatePoint(point, new_rg)
//			}
//			//fmt.Printf("Interval found %+v %+v \n", int_container, interval)
//		} else {
//			//fmt.Println("Empty interval:", int_container, interval)
//		}
//	}
//}

func (s *Subspace) Cluster(min_dense_points int, min_cluster_points int){
	s.Grid.Cluster(min_dense_points, min_cluster_points)
}