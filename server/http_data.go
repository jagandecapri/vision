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

type Color string

type Point_data struct{
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Point_metadata struct{
	Color Color `json:"color"`
}

type Point struct{
	Point_data Point_data `json:"data"`
	Point_metadata Point_metadata `json:"metadata"`
}

type Graph_metadata struct{
	ID string `json:"id"`
}

type Graph struct{
	Graph_metadata Graph_metadata `json:"metadata"`
	Points []Point `json:"points"`
}

type HttpData1 []Graph