package tree

type Cluster struct{
	Cluster_id int
	Cluster_type int
	Num_of_points int
	ListOfUnits map[Range]*Unit
}