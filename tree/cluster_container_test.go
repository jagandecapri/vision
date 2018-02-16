package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClusterContainer_GetCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1}
	c2 := Cluster{Cluster_id: 2}
	c3 := Cluster{Cluster_id: 3}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	c1_tmp, _  := cc.GetCluster(1)
	c2_tmp, _  := cc.GetCluster(2)
	c3_tmp, _  := cc.GetCluster(3)

	assert.Equal(t, c1, c1_tmp)
	assert.Equal(t, c2, c2_tmp)
	assert.Equal(t, c3, c3_tmp)
}

func TestClusterContainer_AddUpdateCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Num_of_points: 10}
	c2 := Cluster{Cluster_id: 2, Num_of_points: 10}
	c3 := Cluster{Cluster_id: 3, Num_of_points: 20}
	
	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	tmp := Cluster{Cluster_id: 3}
	tmp1 := Cluster{Cluster_id: 4}

	cc.AddUpdateCluster(tmp)
	cc.AddUpdateCluster(tmp1)

	assert.Equal(t, tmp, cc.ListOfClusters[3])
	assert.Equal(t, tmp1, cc.ListOfClusters[4])
}

func TestClusterContainer_GetBiggestCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Num_of_points: 10}
	c2 := Cluster{Cluster_id: 2, Num_of_points: 20}
	c3 := Cluster{Cluster_id: 3, Num_of_points: 10}

	cc := ClusterContainer{ListOfClusters: map[int]Cluster{}}

	cc.AddUpdateCluster(c1)
	cc.AddUpdateCluster(c2)
	cc.AddUpdateCluster(c3)

	assert.Equal(t, c2, cc.GetBiggestCluster())

	c4 := Cluster{Cluster_id: 1, Num_of_points: 10}
	c5 := Cluster{Cluster_id: 2, Num_of_points: 10}
	c6 := Cluster{Cluster_id: 3, Num_of_points: 10}

	cc1 := ClusterContainer{ListOfClusters: map[int]Cluster{}}

	cc1.AddUpdateCluster(c4)
	cc1.AddUpdateCluster(c5)
	cc1.AddUpdateCluster(c6)

	assert.Equal(t, c6, cc1.GetBiggestCluster())
}

func TestClusterContainer_RemoveCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1}
	c2 := Cluster{Cluster_id: 2}
	c3 := Cluster{Cluster_id: 3}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	tmp := Cluster{Cluster_id: 3, ListOfUnits: map[Range]*Unit{}}

	removed_cluster := cc.RemoveCluster(3)

	assert.Equal(t, map[int]Cluster{1: c1, 2: c2}, cc.ListOfClusters)
	assert.Equal(t, tmp, removed_cluster)
}