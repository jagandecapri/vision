package tree

type ClusterContainer struct{
	BiggestCluster Cluster
	ListOfClusters map[int]Cluster
}

func (cc *ClusterContainer) GetBiggestCluster() Cluster{
	return cc.BiggestCluster
}

func (cc *ClusterContainer) GetCluster(cluster_id int) (Cluster, bool){
	cluster, ok := cc.ListOfClusters[cluster_id]
	return cluster, ok
}

func (cc *ClusterContainer) GetClusters() map[int]Cluster{
	return cc.ListOfClusters
}

func (cc *ClusterContainer) AddUpdateCluster(cluster Cluster){
	if cluster.Num_of_points >= cc.BiggestCluster.Num_of_points{
		cc.BiggestCluster = cluster
	}
	cc.ListOfClusters[cluster.Cluster_id] = cluster
}

func (cc *ClusterContainer) RemoveCluster(cluster_id int) Cluster{
	cluster_remove := cc.ListOfClusters[cluster_id]
	tmp := map[Range]*Unit{}
	for rg, unit := range cluster_remove.ListOfUnits{
		tmp[rg] = unit
	}
	tmp_cluster := Cluster{Cluster_id: cluster_remove.Cluster_id,
	//Cluster_type: cluster_remove.Cluster_type,
	ListOfUnits: tmp}
	delete(cc.ListOfClusters, cluster_id)
	return tmp_cluster
}