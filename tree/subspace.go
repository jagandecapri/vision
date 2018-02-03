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
		tmp := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		tmp1 := [2]float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
		int_container := IntervalContainer{Id: 1, Range: Range{Low: tmp, High: tmp1}, Scale_factor: s.Scale_factor}
		interval := (*s.interval_tree).Query(int_container)
		if len(interval) > 0{
			interval_ext := interval[0].(IntervalContainer)
			Vec := []float64{p.Vec_map[subspace_key[0]], p.Vec_map[subspace_key[1]]}
			point := Point{
				Unit_id:  int(interval[0].ID()),
				Vec: Vec,
				Id: p.Id,
				Vec_map:  map[string]float64{},
			}

			for k, v := range p.Vec_map{
				point.Vec_map[k] = v
			}

			cur_rg := s.Grid.GetPointRange(point.GetID())
			new_rg := interval_ext.Range
			if cur_rg == (Range{}){
				//fmt.Println("Add point called")
				s.Grid.AddPoint(point, new_rg)
			} else if cur_rg != (Range{}) && cur_rg != new_rg{
				//fmt.Println("Update point called")
				s.Grid.UpdatePoint(point, new_rg)
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