package tree

type KDTree_Extend struct{
	*KDTree
	Units
}

func (kd_ext *KDTree_Extend) AddToStore(key Range, p PointContainer){
	tmp := kd_ext.Units.Store[key]
	tmp.AddPoint(p)
}