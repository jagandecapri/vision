package tree

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T){
	points := []point{{x:1.0, y:3.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:1.0, y:4.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:1.0, y:5.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:1.0, y:6.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:2.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:3.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:4.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:5.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:6.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:2.0, y:7.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:2.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:3.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:4.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:5.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:6.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:3.0, y:7.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:2.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:3.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:4.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:5.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:6.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:7.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:4.0, y:8.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:2.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:3.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:4.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:5.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:6.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:7.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:5.0, y:8.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:2.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:3.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:4.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:5.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:6.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:6.0, y:7.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:7.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:7.0, y:2.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:7.0, y:3.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:7.0, y:4.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:7.0, y:5.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:1.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:2.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:3.0, z:0.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:4.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:5.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:8.0, y:6.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:9.0, y:2.0, z:1.0, cluster_id: UNCLASSIFIED},
		{x:9.0, y:3.0, z:1.0, cluster_id: UNCLASSIFIED}}
	expected_cluster_id := []int{0, 0, 0, 0, 2, 1, 1, 1, 1, 3, 2, 2, 2, 1, 1, 3, 3, 2, 2, 1, 1, -2,
		1, 3, 3, 2, 1, 1, 1, 1, 1, 3, 3, 2, 1, 3, 3, 3, 3, 3, 2, 1, 1, 1, 3, 2, 2, 1, 3, 3, 3, 2, 2}
	res_points := Main_mock(points, 1.0, 2.0)

	for i := 0; i < len(res_points); i++{
		assert.Equal(t, expected_cluster_id[i], res_points[i].cluster_id)
		//TODO: Something wrong with the noise logic
	}
}