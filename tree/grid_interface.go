package tree

type GridInterface interface{
	GetUnits() map[Range]*Unit
	GetMinDensePoints() int
	GetMinClusterPoints() int
	GetNextClusterID() int
}
