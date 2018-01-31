package server

type PointCluster map[int][][]float64

type HttpData map[string]PointCluster
//
//type PointCluster struct{
//	Id int
//	Cluster_id int
//	Cluster_type string
//}
//
//type HttpData struct{
//	Point_cluster map[string]map[int]PointCluster
//	Points []tree.PointContainer
//}