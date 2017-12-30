package tree

import (
	"github.com/jagandecapri/vision/utils"
)

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

func IGDCA(grid Grid, min_dense_points int, min_cluster_points int) (map[Range]*Unit, map[int]Cluster){
	units := grid.GetUnits()
	cluster_id := grid.GetNextClusterID()
	cluster_map := grid.GetClusterMap()

	for rg, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent
				var ret int
				var neighbour_cluster_ids []int
				ret, neighbour_cluster_ids, cluster_map = AbsorbIntoCluster(unit, rg, cluster_map, min_dense_points)
				if ret == SUCCESS{
					if len(neighbour_cluster_ids) > 1{
						_, neighbour_cluster_ids, cluster_map = MergeClusters(cluster_map, neighbour_cluster_ids)
					}
					num_points_cluster := 0
					for _, cluster_id := range neighbour_cluster_ids{
						num_points_cluster = ComputeNumberOfPointsInCluster(cluster_map[cluster_id])
						cluster_map[cluster_id] = ComputeClusterType(min_cluster_points, num_points_cluster, cluster_map[cluster_id])
					}
				}else if ret == FAILURE{
					cluster_map = NewCluster(unit, rg, cluster_id, min_dense_points, min_cluster_points, cluster_map)
					cluster_id = grid.GetNextClusterID()
				}
			}
		}
	}
	return units, cluster_map
}

func ComputeNumberOfPointsInCluster(cluster Cluster) int{
	num_points_cluster := 0
	for _, unit := range cluster.ListOfUnits{
		num_points_cluster += unit.GetNumberOfPoints()
	}
	return num_points_cluster
}
func ComputeClusterType(min_cluster_points int, num_points_cluster int, cluster Cluster) Cluster{
	if num_points_cluster >= min_cluster_points{
		cluster.Cluster_type = NON_OUTLIER_CLUSTER
	} else {
		cluster.Cluster_type = OUTLIER_CLUSTER
	}
	return cluster
}

func NewCluster(unit *Unit, rg Range, cluster_id int, min_dense_points int, min_cluster_points int, cluster_map map[int]Cluster) map[int]Cluster{
	cluster := Cluster{Cluster_id: cluster_id, ListOfUnits: make(map[Range]*Unit)}
	num_points_cluster := expand(unit, rg, cluster_id, min_dense_points, cluster)
	cluster = ComputeClusterType(min_cluster_points, num_points_cluster, cluster)
	cluster_map[cluster_id] = cluster
	return cluster_map
}

func AbsorbIntoCluster(unit *Unit, rg Range, cluster_map map[int]Cluster, min_dense_points int) (int, []int, map[int]Cluster){
	ret_value := FAILURE
	cluster_ids := []int{}
	for _, neighbour_unit := range unit.Neighbour_units{
		if isDenseUnit(neighbour_unit, min_dense_points) && neighbour_unit.Cluster_id != UNCLASSIFIED &&
			neighbour_unit.Cluster_id != NOISE{
			unit.Cluster_id = neighbour_unit.Cluster_id
			cluster_ids = append(cluster_ids, unit.Cluster_id)
		}
	}
	cluster_ids = utils.UniqInt(cluster_ids)
	if len(cluster_ids) > 0{
		tmp := cluster_ids[0]
		unit.Cluster_id = tmp
		cluster_map[tmp].ListOfUnits[rg] = unit
		ret_value = SUCCESS
	}
	return ret_value, cluster_ids, cluster_map
}

func MergeClusters(cluster_map map[int]Cluster, cluster_ids []int) (int, []int, map[int]Cluster){
	ret_value := FAILURE
	cluster_id_merged, cluster_id_to_be_merged := cluster_ids[0], cluster_ids[1:]
	tmp := []int{cluster_id_merged}
	for _, cluster_id := range cluster_id_to_be_merged{
		cluster, ok := cluster_map[cluster_id]
		if ok{
			for _, unit := range cluster.ListOfUnits{
				unit.Cluster_id = cluster_id_merged
			}
			ret_value = SUCCESS
			delete(cluster_map, cluster_id)
		}
	}
	return ret_value, tmp, cluster_map
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