package tree

import "github.com/jagandecapri/vision/orunada/utils"

const UNCLASSIFIED = 0
const NOISE = -1

const OUTLIER_CLUSTER = -3
const NON_OUTLIER_CLUSTER = -4

const SUCCESS = -10
const FAILURE = -11

type ClusterInterface interface{
	GetUnits() map[Range]*Unit
	GetMinDensePoints() int
	GetMinClusterPoints() int
	GetNextClusterID() int
}

type Cluster struct{
	Cluster_id int
	Cluster_type int
	ListOfUnits map[Range]*Unit
}

func GDA(Units Units, min_dense_points int, min_cluster_points int) (map[Range]*Unit, map[int]Cluster){
	units := Units.GetUnits()
	cluster_id := Units.GetNextClusterID()
	cluster_map := Units.GetClusterMap()

	for rg, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent
				var ret bool
				var neighbour_cluster_ids []int
				ret, neighbour_cluster_ids = AbsorbIntoCluster(unit, Units, min_dense_points)
				if ret == SUCCESS{
					if len(neighbour_cluster_ids) > 1{
						MergeClusters(units, neighbour_cluster_ids)
					}
				} else if ret == FAILURE{
					NewCluster(unit, rg, cluster_id, min_dense_points, min_cluster_points, cluster_map)
					cluster_id = Units.GetNextClusterID()
				}
			}
		}
	}
	return units, cluster_map
}

func NewCluster(unit *Unit, rg Range, cluster_id int, min_dense_points int, min_cluster_points int, cluster_map map[int]Cluster){
	cluster := Cluster{Cluster_id: cluster_id, ListOfUnits: make(map[Range]*Unit)}
	num_points_cluster := expand(unit, rg, cluster_id, min_dense_points, cluster)
	if num_points_cluster >= min_cluster_points{
	cluster.Cluster_type = NON_OUTLIER_CLUSTER
	} else {
	cluster.Cluster_type = OUTLIER_CLUSTER
	}
	cluster_map[cluster_id] = cluster
}

func AbsorbIntoCluster(unit *Unit, units Units, min_dense_points int) (int, neighbour_cluster_ids []int){
	ret_value := FAILURE
	cluster_ids := []int{}
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points) && neighbour_unit.Cluster_id != UNCLASSIFIED ||
			neighbour_unit.Cluster_id != NOISE{
			unit.Cluster_id = neighbour_unit.Cluster_id
			cluster_ids = append(cluster_ids, unit.Cluster_id)
		}
	}
	cluster_ids = utils.UniqInt(cluster_ids)
	if len(cluster_ids) > 0{
		tmp := cluster_ids[0]
		unit.Cluster_id = tmp
		ret_value = SUCCESS
	}

	return ret_value, cluster_ids
}

func MergeClusters(units Units, cluster_ids []int) int{
	ret_value := FAILURE
	cluster_map := units.GetClusterMap()
	cluster_id_to_merge := cluster_ids[0]
	for _, cluster_id := range cluster_ids{
		cluster, ok := cluster_map[cluster_id]
		if ok{
			for _, unit := range cluster.ListOfUnits{
				unit.Cluster_id = cluster_id_to_merge
			}
			ret_value = SUCCESS
		}
	}
	return ret_value
}

type Seed struct{
	unit *Unit
	rg Range
}

func expand(unit *Unit, rg Range, cluster_id int, min_dense_points int, cluster Cluster) (int){
	point_count_acc := 0

	seeds := []Seed{}
	for rg, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points){
			if neighbour_unit.Cluster_id == UNCLASSIFIED {
				seed := Seed{unit: neighbour_unit, rg: rg}
				seeds = append(seeds, seed)
			}
		}
	}

	unit.Cluster_id = cluster_id
	point_count_acc += unit.GetNumberOfPoints()
	cluster.ListOfUnits[rg] = unit
	point_count_acc = spread(point_count_acc, seeds, cluster_id, min_dense_points, cluster)
	return point_count_acc
}

func spread(point_count_acc int, seeds []Seed, cluster_id int, min_dense_points int, cluster Cluster) int{
	if len(seeds) == 0{
		return point_count_acc
	}
	var seed Seed
	var unit *Unit
	seed, seeds = seeds[0], seeds[1:]
	unit = seed.unit

	if unit.Cluster_id == UNCLASSIFIED || unit.Cluster_id == NOISE{
		unit.Cluster_id = cluster_id
		point_count_acc += unit.GetNumberOfPoints()
		cluster.ListOfUnits[seed.rg] = unit

		for rg, neighbour_unit := range unit.Neighbour_units{
			if isDenseUnit(neighbour_unit, min_dense_points){
				if neighbour_unit.Cluster_id == UNCLASSIFIED {
					seed := Seed{unit: neighbour_unit, rg: rg}
					seeds = append(seeds, seed)
				}
			}
		}
	}

	return spread(point_count_acc, seeds, cluster_id, min_dense_points, cluster)
}

func isDenseUnit(unit *Unit, min_dense_points int) bool{
	return unit.GetNumberOfPoints() >= min_dense_points
}