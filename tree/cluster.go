package tree

type Cluster struct{
	Cluster_id int
	Cluster_type int
	ListOfUnits map[Range]*Unit
}