package tree

type KDTree_Extend struct{
	*KDTree
	Store map[int]*Unit
}

func (kd_ext *KDTree_Extend) AddUnit(unit *Unit){
	key := unit.GetID()
	kd_ext.Store[key] = unit
	kd_ext.Insert(unit)
}

func (kd_ext *KDTree_Extend) AddToStore(key int, p PointContainer){
	tmp := kd_ext.Store[key]
	tmp.AddPoint(p)
}