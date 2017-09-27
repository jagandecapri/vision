package tree

const UNCLASSIFIED = 0
const NOISE = -1

func Cluster(kd_tree_ext *KDTree_Extend, min_dense_points int, min_cluster_points int){
    out := make(chan PointInterface)
    kd_tree_ext.BFSTraverseChan(out)
    cluster_id := 1
    for point := range out{
        unit := Unit(point)
        if unit.Cluster_id == UNCLASSIFIED {
            if cluster_id == 1{
                unit.Cluster_id = 1
            } else {
                for _, neighbour_unit := range unit.Neighbour_units{
                    if neighbour_unit == UNCLASSIFIED  neighbour_unit != NOISE{
                        unit.Cluster_id = neighbour_unit
                    }
                }
            }
        }
    }
}