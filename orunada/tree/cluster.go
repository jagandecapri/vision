package tree

const UNCLASSIFIED = 0
const NOISE = -1

const CORE_UNIT = 1
const NOT_CORE_UNIT = 2

func GDA(units Units, min_dense_points int, min_cluster_points int) (Units, []int){
	cluster_id := 1
	non_outlier_clusters := []int{}
	for _, unit := range units.Store{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){
				unit_type, num_points_cluster := expand(unit, cluster_id, min_dense_points)
				if num_points_cluster >= min_cluster_points{
					non_outlier_clusters = append(non_outlier_clusters, cluster_id)
				}
				//else {
				//	fmt.Println("Cluster id: ", cluster_id, "num_points_cluster: ", num_points_cluster, )
				//	for _, unit := range units.Store{
				//		if unit.Cluster_id == cluster_id{
				//			fmt.Println("Unit point count: ", unit.GetNumberOfPoints())
				//		}
				//	}
				//}
				if (unit_type == CORE_UNIT){
					cluster_id++
				}
			}
		}
	}
	return units, non_outlier_clusters
}

func expand(unit *Unit, cluster_id int, min_dense_points int) (int, int){
	return_value := NOT_CORE_UNIT
	point_count_acc := 0

	seeds := []*Unit{}
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points){
			seeds = append(seeds, neighbour_unit)
		}
	}

	if len(seeds) == 0{
		unit.Cluster_id = NOISE
	} else {
		unit.Cluster_id = cluster_id
		point_count_acc += unit.GetNumberOfPoints()
		for _, neighbour_unit := range seeds{
			neighbour_unit.Cluster_id = cluster_id
			point_count_acc += unit.GetNumberOfPoints()
		}
		for _, neighbour_unit := range seeds{
			point_count_acc += spread(neighbour_unit, seeds, cluster_id, min_dense_points)
		}
		return_value = CORE_UNIT
	}
	return return_value, point_count_acc
}

func spread(unit *Unit, seeds []*Unit, cluster_id int, min_dense_points int) int{
	spread := []*Unit{}
	point_count_acc := 0
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points){
			spread = append(spread, neighbour_unit)
		}
	}
	if len(spread) > 0{
		for _, neighbour_unit := range spread{
			if neighbour_unit.Cluster_id == UNCLASSIFIED || neighbour_unit.Cluster_id == NOISE{
				if neighbour_unit.Cluster_id == UNCLASSIFIED{
					seeds = append(seeds, neighbour_unit)
				}
				neighbour_unit.Cluster_id = cluster_id
				point_count_acc += unit.GetNumberOfPoints()
			}
		}
	}
	return point_count_acc
}

func isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}