package tree

const FIRST_CLUSTER = 1

func Cluster(kd_tree_ext *KDTree_Extend, min_dense_points int, min_cluster_points int){
    out := make(chan PointInterface)
    kd_tree_ext.BFSTraverseChan(out)
    cluster_id = FIRST_CLUSTER
    for point := range out{
        unit := Unit(point)
        if unit.Cluster_id == UNCLUSTERED {
            if cluster_id == FIRST_CLUSTER{
                unit.Cluster_id = FIRST_CLUSTER
            } else {
                for _, neighbour_unit := unit.Neighbour_units{
                    if neighbour_unit != UNCLUSTERED || neighbour_unit != NOISE{
                        unit.Cluster_id = neighbour_unit
                    }
                }
            }
        }
    }
}