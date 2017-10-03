package tree

const UNCLASSIFIED = 0
const NOISE = -1

const CORE_UNIT = 1
const NOT_CORE_UNIT = 2

var cluster_map map[int]Unit

func Cluster(units Units, min_dense_points int, min_cluster_points int) {
	cluster_id := 1
	for rg, unit := range units.Store {
		seeds := Queue{}
		if isDenseUnit(unit, min_dense_points) == true {
			unit.Cluster_id = cluster_id
			for _, neighbour_unit := range unit.Neighbour_units {
				if neighbour_unit.Cluster_id == UNCLASSIFIED || neighbour_unit.Cluster_id == NOISE {
					if neighbour_unit.GetNumberOfPoints() >= min_dense_points {
						neighbour_unit.Cluster_id = cluster_id
						seeds.Push(neighbour_unit)
					}
				}
			}
			if len(seeds) == 0 {
				unit.Cluster_id = NOISE
			} else {
				for seed := seeds.Pop(); seed != nil; {
					expandCluster(seed, cluster_id, min_dense_points, seeds)
				}
			}
		}
		cluster_id++
                units.Store[rg] = unit
	}
}

func expandCluster(unit *Unit, cluster_id int, min_dense_points int, seeds Queue) {
	for _, neighbour_unit := range unit.Neighbour_units {
		if neighbour_unit.Cluster_id == UNCLASSIFIED || neighbour_unit.Cluster_id == NOISE {
			unit.Cluster_id = cluster_id
		}
		if isDenseUnit(neighbour_unit, min_dense_points) {
			seeds.Push(neighbour_unit)
		}
	}
}

func isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}