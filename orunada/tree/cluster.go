package tree

const UNCLASSIFIED = 0
const NOISE = -1

func GDA(units map[Range]*Unit, min_dense_points int, min_cluster_points int) (map[Range]*Unit, map[string][]int){
	cluster_id := 1
	cluster_map := map[string][]int{
		"outlier_clusters": []int{},
		"non_outlier_clusters": []int{},
	}
	for _, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent
				num_points_cluster := expand(unit, cluster_id, min_dense_points)
				if num_points_cluster >= min_cluster_points{
					cluster_map["non_outlier_clusters"] = append(cluster_map["non_outlier_clusters"], cluster_id)
				} else {
					cluster_map["outlier_clusters"] = append(cluster_map["outlier_clusters"], cluster_id)
				}
				cluster_id++
			}
		}
	}
	return units, cluster_map
}

func expand(unit *Unit, cluster_id int, min_dense_points int) (int){
	point_count_acc := 0

	seeds := []*Unit{}
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points){
			seeds = append(seeds, neighbour_unit)
		}
	}

	unit.Cluster_id = cluster_id
	point_count_acc += unit.GetNumberOfPoints()
	for _, neighbour_unit := range seeds{
		neighbour_unit.Cluster_id = cluster_id
		point_count_acc += unit.GetNumberOfPoints()
	}
	for _, neighbour_unit := range seeds{
		point_count_acc += spread(neighbour_unit, seeds, cluster_id, min_dense_points)
	}

	return point_count_acc
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