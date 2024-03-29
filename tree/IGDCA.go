package tree

import (
	"github.com/jagandecapri/vision/utils"
	"log"
)

const UNCLASSIFIED = 0
const NOISE = -1

const SUCCESS = -10
const FAILURE = -11

func IGDCA(grid *Grid, min_dense_points int, min_cluster_points int) (map[Range]*Unit){
	units := grid.GetUnitsToCluster()
	cluster_id := grid.GetNextClusterID()

	for rg, unit := range units{
		if unit.Cluster_id == UNCLASSIFIED{
			if isDenseUnit(unit, min_dense_points){ //TODO: Could be redundant is only dense units are sent

				ret, cluster, cluster_ids_to_be_merged := AbsorbIntoCluster(grid, unit, min_dense_points)
				if ret == SUCCESS{
					if len(cluster_ids_to_be_merged) > 1{
						_, cluster = MergeClusters(grid, cluster, cluster_ids_to_be_merged)
					}

					if isClusterTooSmall(min_cluster_points, cluster) == true{
						log.Println("Cluster to remove: ", cluster.Cluster_id)
						grid.RemoveCluster(cluster.Cluster_id)
					} else {
						grid.AddUpdateCluster(cluster)
					}
				}else if ret == FAILURE{
					cluster := NewCluster(unit, rg, cluster_id, min_dense_points)
					grid.AddUpdateCluster(cluster)
					cluster_id = grid.GetNextClusterID()
				}
			}
		}
	}
	return units
}

func isClusterTooSmall(min_cluster_points int, cluster Cluster) bool{
	return cluster.Num_of_points < min_cluster_points
}

func NewCluster(unit *Unit, rg Range, cluster_id int, min_dense_points int) (Cluster){
	cluster := Cluster{Cluster_id: cluster_id, ListOfUnits: make(map[Range]*Unit)}
	cluster = expand(unit, rg, cluster_id, min_dense_points, cluster)
	return cluster
}

func AbsorbIntoCluster(grid *Grid, unit *Unit, min_dense_points int) (int, Cluster, []int){
	ret_value := FAILURE
	var cluster_ids []int
	var cluster Cluster
	var ok bool

	for _, neighbour_unit := range unit.GetNeighbouringUnits() {
		if isDenseUnit(neighbour_unit, min_dense_points) && neighbour_unit.Cluster_id != UNCLASSIFIED &&
			neighbour_unit.Cluster_id != NOISE{
			cluster_ids = append(cluster_ids, neighbour_unit.Cluster_id)
		}
	}

	cluster_ids = utils.UniqInt(cluster_ids)
	var unit_cluster_id int
	var cluster_id_to_be_merged []int

	if len(cluster_ids) > 0{
		unit_cluster_id = cluster_ids[0]
		cluster_id_to_be_merged = cluster_ids[1:]
		ret_value = SUCCESS

		unit.Cluster_id = unit_cluster_id
		cluster, ok = grid.GetCluster(unit_cluster_id) //ok value here need to be taken into account if want to avoid nil panic
		if ok{
			cluster.ListOfUnits[unit.Range] = unit
			cluster.Num_of_points += unit.GetNumberOfPoints()
		} else {
			clusters := grid.GetClusters()
			log.Printf("Cluster id requested: %v Cluster ID available: ", unit_cluster_id)
			for _, cluster := range clusters{
				log.Printf("%v ", cluster.Cluster_id)
			}
			for _, neighbour_unit := range unit.GetNeighbouringUnits() {
				if isDenseUnit(neighbour_unit, min_dense_points) && neighbour_unit.Cluster_id == unit_cluster_id{
					log.Printf(" Rogue neighbour unit: %v ", neighbour_unit.Id)
				}
			}
			log.Printf("\n")
		}
	}
	return ret_value, cluster, cluster_id_to_be_merged
}

func MergeClusters(grid *Grid, cluster Cluster, cluster_ids_to_be_merged []int) (int, Cluster){
	ret_value := FAILURE

	for _, cluster_id := range cluster_ids_to_be_merged{
		cluster_to_merge, ok := grid.GetCluster(cluster_id)
		if ok{
			for rg, unit := range cluster_to_merge.ListOfUnits{
				unit.Cluster_id = cluster.Cluster_id
				cluster.ListOfUnits[rg] = unit
				cluster.Num_of_points += unit.GetNumberOfPoints()
			}
			ret_value = SUCCESS
			log.Println("Cluster to remove after merge: ", cluster_id)
			grid.RemoveCluster(cluster_id)
		}
	}

	return ret_value, cluster
}

type Seed struct{
	unit *Unit
	rg Range
}

func expand(unit *Unit, rg Range, cluster_id int, min_dense_points int, cluster Cluster) (Cluster){
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
	cluster.Num_of_points += unit.GetNumberOfPoints()
	cluster = spread(point_count_acc, seeds, cluster_id, min_dense_points, cluster)
	return cluster
}

func spread(point_count_acc int, seeds []Seed, cluster_id int, min_dense_points int, cluster Cluster) (Cluster){
	if len(seeds) == 0{
		return cluster
	}
	var seed Seed
	var unit *Unit
	seed, seeds = seeds[0], seeds[1:]
	unit = seed.unit

	if unit.Cluster_id == UNCLASSIFIED || unit.Cluster_id == NOISE{
		unit.Cluster_id = cluster_id
		cluster.ListOfUnits[seed.rg] = unit
		cluster.Num_of_points += unit.GetNumberOfPoints()

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