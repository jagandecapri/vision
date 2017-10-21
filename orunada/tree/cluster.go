package tree

const UNCLASSIFIED = 0
const NOISE = -1

const OUTLIER_CLUSTER = -3
const NON_OUTLIER_CLUSTER = -4

type ClusterInterface interface{
	GetUnits() map[Range]*Unit
	GetMinDensePoints() int
	GetMinClusterPoints() int
	GetNextClusterID() int
}

type Cluster struct{
	Cluster_id int
	Cluster_type int
	ListOfUnits []*Unit
}

func GDA(units map[Range]*Unit, min_dense_points int, min_cluster_points int) (map[Range]*Unit, map[int]Cluster){
	cluster_id := 1
	cluster_map := make(map[int]Cluster)

	for _, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent
				cluster := Cluster{Cluster_id: cluster_id}
				num_points_cluster := expand(unit, cluster_id, min_dense_points, cluster)
				if num_points_cluster >= min_cluster_points{
					cluster.Cluster_type = NON_OUTLIER_CLUSTER
				} else {
					cluster.Cluster_type = OUTLIER_CLUSTER
				}
				cluster_map[cluster_id] = cluster
				cluster_id++
			}
		}
	}
	return units, cluster_map
}

func expand(unit *Unit, cluster_id int, min_dense_points int, cluster Cluster) (int){
	point_count_acc := 0

	seeds := []*Unit{}
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points){
			if neighbour_unit.Cluster_id == UNCLASSIFIED {
				seeds = append(seeds, neighbour_unit)
			}
		}
	}

	unit.Cluster_id = cluster_id
	point_count_acc += unit.GetNumberOfPoints()
	cluster.ListOfUnits = append(cluster.ListOfUnits, unit)
	point_count_acc = spread(point_count_acc, seeds, cluster_id, min_dense_points, cluster)
	return point_count_acc
}

func spread(point_count_acc int, seeds []*Unit, cluster_id int, min_dense_points int, cluster Cluster) int{
	if len(seeds) == 0{
		return point_count_acc
	}
	var unit *Unit
	unit, seeds = seeds[0], seeds[1:]

	if unit.Cluster_id == UNCLASSIFIED || unit.Cluster_id == NOISE{
		unit.Cluster_id = cluster_id
		point_count_acc += unit.GetNumberOfPoints()
		cluster.ListOfUnits = append(cluster.ListOfUnits, unit)

		for _, neighbour_unit := range unit.Neighbour_units{
			if isDenseUnit(neighbour_unit, min_dense_points){
				if neighbour_unit.Cluster_id == UNCLASSIFIED {
					seeds = append(seeds, neighbour_unit)
				}
			}
		}
	}

	return spread(point_count_acc, seeds, cluster_id, min_dense_points, cluster)
}

func isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}