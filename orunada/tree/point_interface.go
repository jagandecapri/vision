package tree

type PointInterface interface {
	// Return the total number of dimensions
	Dim() int
	// Return the value X_{dim}, dim is started from 0
	GetValue(dim int) int
	// Return the distance between two points
	Distance(point PointInterface) float64
	// Return the distance between the point and the plane X_{dim}=val
	PlaneDistance(val float64, dim int) float64
	// ID
	GetID() int
}
