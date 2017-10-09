package tree


type Units struct{
	Store map[Range]*Unit
	listDenseUnits map[Range]*Unit
	listOldDenseUnits map[Range]*Unit
	listNewDenseUnits map[Range]*Unit
	listUnitToRep map[Range]*Unit
}


func (us Units) AddUnit(unit *Unit, rg Range){
	us.Store[rg] = unit
}

func (us Units) SetupGrid(interval_l float64){
	for rg, unit := range us.Store{
		unit.Neighbour_units = us.GetNeighbouringUnits(rg, interval_l)
		us.Store[rg] = unit
	}
}

func (us Units) RecomputeDenseUnits(min_dense_points int){
	for rg, unit := range us.Store{
		if isDenseUnit(unit, min_dense_points){
			_, ok := us.listDenseUnits[rg]
			if !ok{
				us.listDenseUnits[rg] = unit
				us.listNewDenseUnits[rg] = unit
			}
		} else {
			_, ok := us.listDenseUnits[rg]
			if ok{
				delete(us.listDenseUnits, rg)
				us.listOldDenseUnits[rg] = unit
			}
		}
	}
	us.ProcessOldDenseUnits()
}

func (us Units) ProcessOldDenseUnits(){
	for _, unit := range us.listOldDenseUnits{
		cluster_id := unit.Cluster_id
		unit.Cluster_id = UNCLASSIFIED
		count_neighbour_same_cluster := 0

		for _, neighbour_unit := range unit.Neighbour_units{
			if neighbour_unit.Cluster_id == cluster_id{
				count_neighbour_same_cluster++
			}
			if count_neighbour_same_cluster > 2{
				break
			}
		}
		if count_neighbour_same_cluster > 2 {
			//TODO: Write for-loop to unclusterise all units having the cluster id and add to listOrRep
		}
	}
}

func (us Units) isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}

func (us Units) GetNeighbouringUnits(rg Range, interval_l float64) []*Unit {
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

	neighbour_units := []*Unit{}

	for _, rg := range tmp{
		if unit, ok := us.Store[rg]; ok{
			neighbour_units = append(neighbour_units, unit)
		}
	}

	return neighbour_units
}