package tree

type Grid struct{
	Store              map[Range]*Unit
	ClusterContainer
	point_unit_map     map[int]Range
	minDensePoints     int
	minClusterPoints   int
	cluster_id_counter int
	listDenseUnits     map[Range]*Unit
	tmpUnitToCluster   map[Range]*Unit
}

func NewGrid() Grid {
	units := Grid{
		Store: make(map[Range]*Unit),
		point_unit_map: make(map[int]Range),
		listDenseUnits: make(map[Range]*Unit),
		tmpUnitToCluster: make(map[Range]*Unit),
		ClusterContainer: ClusterContainer{ListOfClusters: make(map[int]Cluster)},
	}
	return units
}

func (us *Grid) GetUnits() map[Range]*Unit{
	return us.Store
}

func (us *Grid) GetUnitsToCluster() map[Range]*Unit{
	if len(us.tmpUnitToCluster) == 0{
		return us.Store
	} else {
		return us.tmpUnitToCluster
	}
}

func (us *Grid) GetMinDensePoints() int{
	return us.minDensePoints
}

func (us *Grid) GetMinClusterPoints() int{
	return us.minClusterPoints
}

func (us *Grid) GetNextClusterID() int{
	us.cluster_id_counter += 1
	return us.cluster_id_counter
}

func (us *Grid) RemovePoint(point Point, rg Range){
	unit, ok := us.Store[rg]
	if ok{
		unit.RemovePoint(point)
		delete(us.point_unit_map, point.GetID())
	}
}

func (us *Grid) AddPoint(point Point, rg Range){
	unit, ok := us.Store[rg]
	if ok{
		unit.AddPoint(point)
		us.point_unit_map[point.GetID()] = rg
	}
}

func (us *Grid) UpdatePoint(point Point, new_range Range){
	point_id := point.GetID()
	cur_range := us.point_unit_map[point_id]
	if cur_range != new_range{
		us.RemovePoint(point, cur_range)
		us.AddPoint(point, new_range)
	}
}

func (us *Grid) GetPointRange(id int) Range{
	return us.point_unit_map[id]
}

func (us *Grid) AddUnit(unit *Unit){
	us.Store[unit.Range] = unit
}

//func (us *Grid) AddUnit(unit *Unit, rg Range){
//	us.Store[rg] = unit
//}

func (us *Grid) SetupGrid(interval_l float64){
	for rg, unit := range us.Store{
		neighbour_units := us.GetNeighbouringUnits(rg, interval_l)
		unit.SetNeighbouringUnits(neighbour_units)
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

		for _, neighbour_unit := range unit.GetNeighbouringUnits() {
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

func (us *Grid) Cluster(min_dense_points int, min_cluster_points int){
	listNewDenseUnits, listOldDenseUnits := us.RecomputeDenseUnits(min_dense_points)
	us.tmpUnitToCluster = listNewDenseUnits
	_ = IGDCA(*us, min_dense_points, min_cluster_points)

	listUnitToRep := us.ProcessOldDenseUnits(listOldDenseUnits)
	us.tmpUnitToCluster = listUnitToRep
	_ = IGDCA(*us, min_dense_points, min_cluster_points)
}

func (us *Grid) GetOutliers() []Point{
	tmp := []Point{}
	for _, unit := range us.Store{
		if unit.Cluster_id == UNCLASSIFIED{
			for _, point := range unit.GetPoints(){
				tmp = append(tmp, point)
			}
		}
	}
	return tmp
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