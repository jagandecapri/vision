package tree

type Grid struct{
	Store map[Range]*Unit
	Point_unit_map map[int]Range
	Cluster_map map[int]Cluster
	MinDensePoints int
	MinClusterPoints int
	cluster_id_counter int
	listDenseUnits map[Range]*Unit
	tmpUnitToCluster map[Range]*Unit
	ClusterContainer
}

func NewGrid() Grid {
	units := Grid{
		Store: make(map[Range]*Unit),
		Point_unit_map: make(map[int]Range),
		Cluster_map: make(map[int]Cluster),
		listDenseUnits: make(map[Range]*Unit),
		tmpUnitToCluster: make(map[Range]*Unit),
		ClusterContainer: ClusterContainer{ListOfClusters: make(map[int]Cluster)},
	}
	return units
}

func (us *Grid) GetUnits() map[Range]*Unit{
	if len(us.tmpUnitToCluster) == 0{
		return us.Store
	} else {
		return us.tmpUnitToCluster
	}
}

func (us *Grid) GetMinDensePoints() int{
	return us.MinDensePoints
}

func (us *Grid) GetMinClusterPoints() int{
	return us.MinClusterPoints
}

func (us *Grid) GetNextClusterID() int{
	us.cluster_id_counter += 1
	return us.cluster_id_counter
}

func (us *Grid) GetClusterMap() map[int]Cluster{
	return us.Cluster_map
}

func (us *Grid) RemovePoint(point PointContainer, rg Range){
	unit, ok := us.Store[rg]
	if ok{
		unit.RemovePoint(point)
		delete(us.Point_unit_map, point.GetID())
	}
}

func (us *Grid) AddPoint(point PointContainer, rg Range){
	unit, ok := us.Store[rg]
	if ok{
		unit.AddPoint(point)
		us.Point_unit_map[point.GetID()] = rg
	}
}

func (us *Grid) UpdatePoint(point PointContainer, new_range Range){
	point_id := point.GetID()
	cur_range := us.Point_unit_map[point_id]
	if cur_range != new_range{
		us.RemovePoint(point, cur_range)
		us.AddPoint(point, new_range)
	}
}

func (us *Grid) GetPointRange(id int) Range{
	return us.Point_unit_map[id]
}

func (us *Grid) AddUnit(unit *Unit, rg Range){
	us.Store[rg] = unit
}

func (us *Grid) SetupGrid(interval_l float64){
	for rg, unit := range us.Store{
		unit.Neighbour_units = us.GetNeighbouringUnits(rg, interval_l)
		us.Store[rg] = unit
	}
}

func (us *Grid) RecomputeDenseUnits(min_dense_points int) (map[Range]*Unit, map[Range]*Unit){
	listNewDenseUnits := make(map[Range]*Unit)
	listOldDenseUnits := make(map[Range]*Unit)

	for rg, unit := range us.Store{
		if isDenseUnit(unit, min_dense_points){
			_, ok := us.listDenseUnits[rg]
			if !ok{
				us.listDenseUnits[rg] = unit
				listNewDenseUnits[rg] = unit
			}
		} else {
			_, ok := us.listDenseUnits[rg]
			if ok{
				delete(us.listDenseUnits, rg)
				listOldDenseUnits[rg] = unit
			}
		}
	}
	return listNewDenseUnits, listOldDenseUnits
}

func (us *Grid) ProcessOldDenseUnits(listOldDenseUnits map[Range]*Unit) map[Range]*Unit {
	listUnitToRep := make(map[Range]*Unit)
	for _, unit := range listOldDenseUnits{
		cluster_id := unit.Cluster_id
		unit.Cluster_id = UNCLASSIFIED
		count_neighbour_same_cluster := 0

		for _, neighbour_unit := range unit.Neighbour_units{
			if neighbour_unit.Cluster_id == cluster_id{
				count_neighbour_same_cluster++
			}
			if count_neighbour_same_cluster >= 2{
				break
			}
		}

		if count_neighbour_same_cluster >= 2 {
			src := us.RemoveCluster(cluster_id)
			for rg, unit := range src.ListOfUnits{
				listUnitToRep[rg] = unit
			}
		}
	}
	return listUnitToRep
}

//func (us *Grid) RemoveCluster(cluster_id int) map[Range]*Unit{
//	tmp := make(map[Range]*Unit)
//	for rg, unit := range us.Cluster_map[cluster_id].ListOfUnits{
//		tmp[rg] = unit
//	}
//	delete(us.Cluster_map, cluster_id)
//	return tmp
//}

func (us *Grid) Cluster(min_dense_points int, min_cluster_points int){
	listNewDenseUnits, listOldDenseUnits := us.RecomputeDenseUnits(min_dense_points)
	us.tmpUnitToCluster = listNewDenseUnits
	_ = IGDCA(*us, min_dense_points, min_cluster_points)

	listUnitToRep := us.ProcessOldDenseUnits(listOldDenseUnits)
	us.tmpUnitToCluster = listUnitToRep
	_ = IGDCA(*us, min_dense_points, min_cluster_points)
}

func (us *Grid) isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}

func (us *Grid) GetNeighbouringUnits(rg Range, interval_l float64) map[Range]*Unit {
	/**
	U = unit; n{x} => neighbouring units
	|n3|n5|n8|
	|n2|U |n7|
	|n1|n4|n6|
	 */
	n1 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.Low[1] - interval_l},
		High: [2]float64{rg.Low[0], rg.Low[1]}}
	n2 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.Low[1]},
		High: [2]float64{rg.Low[0], rg.High[1]}}
	n3 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.High[1]},
		High: [2]float64{rg.Low[0], rg.High[1] + interval_l}}

	n4 := Range{Low:[2]float64{rg.Low[0], rg.Low[1] - interval_l},
		High: [2]float64{rg.High[0], rg.Low[1]}}
	n5 := Range{Low:[2]float64{rg.Low[0], rg.High[1]},
		High: [2]float64{rg.High[0], rg.High[1] + interval_l}}

	n6 := Range{Low:[2]float64{rg.High[0], rg.Low[1] - interval_l},
		High: [2]float64{rg.High[0] + interval_l, rg.Low[1]}}
	n7 := Range{Low:[2]float64{rg.High[0], rg.Low[1]},
		High: [2]float64{rg.High[0] + interval_l, rg.High[1]}}
	n8 := Range{Low:[2]float64{rg.High[0], rg.High[1]},
		High: [2]float64{rg.High[0] + interval_l, rg.High[1] + interval_l}}

	tmp := [8]Range{n1,n2,n3,n4,n5,n6,n7,n8}

	neighbour_units := make(map[Range]*Unit)

	for _, rg := range tmp{
		if unit, ok := us.Store[rg]; ok{
			//neighbour_units = append(neighbour_units, unit)
			neighbour_units[rg] = unit
		}
	}

	return neighbour_units
}