package tree

type KDTree_Extend struct{
	*KDTree
	Store map[int][]IntervalConc
}

func (kd_ext *KDTree_Extend) Add(key int, itv IntervalConc){
	kd_ext.Store[key] = append( kd_ext.Store[key], itv)
}

func (kd_ext *KDTree_Extend) Insert(interval IntervalConc){
	tmp := Point(interval)
	kd_ext.KDTree.Insert(tmp)
}