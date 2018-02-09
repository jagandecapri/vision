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

func IGDCA(grid Grid, min_dense_points int, min_cluster_points int) (map[Range]*Unit){
	units := grid.GetUnitsToCluster()
	cluster_id := grid.GetNextClusterID()

	for rg, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent
				var ret int
				var neighbour_cluster_ids []int
				ret, neighbour_cluster_ids = AbsorbIntoCluster(grid, unit, min_dense_points)
				if ret == SUCCESS{
					if len(neighbour_cluster_ids) > 1{
						_, _, neighbour_cluster_ids = MergeClusters(grid, neighbour_cluster_ids)
					}
					num_points_cluster := 0
					for _, cluster_id := range neighbour_cluster_ids{
						cluster, _ := grid.GetCluster(cluster_id)
						num_points_cluster = ComputeNumberOfPointsInCluster(cluster) //TODO: Optimization to cumulate num_points_cluster in for-loop
						tmp := ComputeClusterType(min_cluster_points, num_points_cluster, cluster)
						grid.AddUpdateCluster(tmp)
					}
				}else if ret == FAILURE{
					num_points_cluster, cluster := NewCluster(unit, rg, cluster_id, min_dense_points)
					cluster = ComputeClusterType(min_cluster_points, num_points_cluster, cluster)
					grid.AddUpdateCluster(cluster)
					cluster_id = grid.GetNextClusterID()
				}
			}
		}
	}
	return units
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

func NewCluster(unit *Unit, rg Range, cluster_id int, min_dense_points int) (int, Cluster){
	cluster := Cluster{Cluster_id: cluster_id, ListOfUnits: make(map[Range]*Unit)}
	num_points_cluster := expand(unit, rg, cluster_id, min_dense_points, cluster)
	return num_points_cluster, cluster
}

func AbsorbIntoCluster(grid Grid, unit *Unit, min_dense_points int) (int, []int){
	ret_value := FAILURE
	cluster_ids := []int{}
	for _, neighbour_unit := range unit.GetNeighbouringUnits() {
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

		cluster, _ := grid.GetCluster(tmp) //ok value here need to tbe taken into account if want to avoid nil panic
		cluster.ListOfUnits[unit.Range] = unit
		grid.AddUpdateCluster(cluster)
		ret_value = SUCCESS
	}
	return ret_value, cluster_ids
}

func MergeClusters(grid Grid, cluster_ids []int) (int, int, []int){
	num_of_points := 0
	ret_value := FAILURE
	cluster_id_merged, cluster_id_to_be_merged := cluster_ids[0], cluster_ids[1:]
	tmp := []int{cluster_id_merged}
	for _, cluster_id := range cluster_id_to_be_merged{
		cluster, ok := grid.GetCluster(cluster_id)
		if ok{
			for _, unit := range cluster.ListOfUnits{
				unit.Cluster_id = cluster_id_merged
				num_of_points += unit.GetNumberOfPoints()
			}
			ret_value = SUCCESS
			grid.RemoveCluster(cluster_id)
		}
	}
	return ret_value, num_of_points, tmp
}

type Seed struct{
	unit *Unit
	rg Range
}

func expand(unit *Unit, rg Range, cluster_id int, min_dense_points int, cluster Cluster) (int){
	point_count_acc := 0

	seeds := []Seed{}
	for rg, neighbour_unit := range unit.GetNeighbouringUnits() {
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

		for rg, neighbour_unit := range unit.GetNeighbouringUnits() {
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