package server

import "github.com/jagandecapri/vision/tree"

type PointCluster struct{
	Id int
	Cluster_id int
	Cluster_type string
}

type HttpData struct{
	Point_cluster map[string]map[int]PointCluster
	Points []tree.PointContainer
}