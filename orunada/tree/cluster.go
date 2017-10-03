package tree

const UNCLASSIFIED = 0
const NOISE = -1

const CORE_UNIT = 1
const NOT_CORE_UNIT = 2

func Cluster(kd_tree_ext *KDTree_Extend, min_dense_points int, min_cluster_points int) {
	out := make(chan PointInterface)
	kd_tree_ext.BFSTraverseChan(out)
	cluster_id := 1
	for point := range out {
		unit := Unit(point)
		seeds := Queue{}
		if isDenseUnit(unit, min_dense_points) {
			unit.Cluster_id = cluster_id
			for _, neighbour_unit := range unit.Neighbour_units {
				if neighbour_unit.Cluster_id == UNCLASSIFIED || neighbour_unit.Cluster_id == NOISE {
					if neighbour_unit.GetNumberOfPoints() >= min_dense_points {
						neighbour_unit.Cluster_id = cluster_id
						seeds.push(neighbour_unit)
					}
				}
			}
			if len(seeds) == 0 {
				unit.Cluster_id == NOISE
			} else {
				for seed := seeds.Pop(); seed != nil; {
					expandCluster(seed, cluster_id, min_dense_points, seeds)
				}
			}
		}
		cluster_id++
	}
}

func expandCluster(unit *Unit, cluster_id int, min_dense_points int, seeds Queue) {
	for _, neighbour_unit := range unit.Neighbour_units {
		if neighbour_unit.Cluster_id == UNCLASSIFIED || neighbour_unit.Cluster_id == NOISE {
			unit.Cluster_id = cluster_id
		}
		if isDenseUnit(neighbour_unit, min_dense_points) {
			seeds.push(neighbour_unit)
		}
	}
}

func isDenseUnit(unit *Unit, min_dense_points int) {
	return unit.GetNumberOfPoints() >= min_dense_points
}
