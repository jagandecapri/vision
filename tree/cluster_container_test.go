package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClusterContainer_GetCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Cluster_type: OUTLIER_CLUSTER}
	c2 := Cluster{Cluster_id: 2, Cluster_type: OUTLIER_CLUSTER}
	c3 := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	c1_tmp, _  := cc.GetCluster(1)
	c2_tmp, _  := cc.GetCluster(2)
	c3_tmp, _  := cc.GetCluster(3)

	assert.Equal(t, c1, c1_tmp)
	assert.Equal(t, c2, c2_tmp)
	assert.Equal(t, c3, c3_tmp)
}

func TestClusterContainer_GetOutliers(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Cluster_type: OUTLIER_CLUSTER}
	c2 := Cluster{Cluster_id: 2, Cluster_type: OUTLIER_CLUSTER}
	c3 := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	assert.Equal(t, map[int]Cluster{1: c1, 2: c2}, cc.GetOutliers())
}

func TestClusterContainer_GetNonOutliers(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Cluster_type: OUTLIER_CLUSTER}
	c2 := Cluster{Cluster_id: 2, Cluster_type: OUTLIER_CLUSTER}
	c3 := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	assert.Equal(t, map[int]Cluster{3: c3}, cc.GetNonOutliers())
}

func TestClusterContainer_AddUpdateCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Cluster_type: OUTLIER_CLUSTER}
	c2 := Cluster{Cluster_id: 2, Cluster_type: OUTLIER_CLUSTER}
	c3 := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	tmp := Cluster{Cluster_id: 3, Cluster_type: OUTLIER_CLUSTER}

	cc.AddUpdateCluster(tmp)

	assert.Equal(t, tmp, cc.ListOfClusters[3])

	tmp1 := Cluster{Cluster_id: 4, Cluster_type: OUTLIER_CLUSTER}

	cc.AddUpdateCluster(tmp1)

	assert.Equal(t, tmp1, cc.ListOfClusters[4])
}

func TestClusterContainer_RemoveCluster(t *testing.T) {
	c1 := Cluster{Cluster_id: 1, Cluster_type: OUTLIER_CLUSTER}
	c2 := Cluster{Cluster_id: 2, Cluster_type: OUTLIER_CLUSTER}
	c3 := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER}


	cc := ClusterContainer{ListOfClusters: map[int]Cluster{1: c1, 2: c2, 3: c3}}

	tmp := Cluster{Cluster_id: 3, Cluster_type: NON_OUTLIER_CLUSTER, ListOfUnits: map[Range]*Unit{}}

	removed_cluster := cc.RemoveCluster(3)

	assert.Equal(t, map[int]Cluster{1: c1, 2: c2}, cc.ListOfClusters)
	assert.Equal(t, tmp, removed_cluster)
}