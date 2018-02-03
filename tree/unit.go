package tree

type Unit struct {
	Id                int
	Cluster_id        int
	Dimension         int
	Center            Point
	Points            map[int]Point
	Center_calculated bool
	Range
	neighbour_units   map[Range]*Unit
}

func (u *Unit) GetNeighbouringUnits() map[Range]*Unit{
	return u.neighbour_units
}

func (u *Unit) SetNeighbouringUnits(neighbour_units map[Range]*Unit){
	u.neighbour_units = neighbour_units
}

func (u *Unit) AddPoint(p Point) {
	u.Center_calculated = false
	u.Points[p.GetID()] = p
}

func (u *Unit) RemovePoint(p Point) {
	u.Center_calculated = false
	delete(u.Points, p.GetID())
}

//Calling GetPoints will update the Cluster_id in each point
func (u *Unit) GetPoints(){

}

func (u *Unit) CalculateCenter() {
	u.GetCenter()
}

func (u *Unit) GetCenter() Point {
	if len(u.Center.Vec) > 0{
		return u.Center
	}
	Center_vec := make([]float64, u.Dimension)
	for _, p := range u.Points {
		for i := 0; i < p.Dim(); i++ {
			Center_vec[i] = Center_vec[i] + p.GetValue(i)
		}
	}
	for i, _ := range Center_vec {
		Center_vec[i] = Center_vec[i] / float64(len(u.Points))
	}
	u.Center = Point{Unit_id: u.Id, Vec: Center_vec}
	return u.Center
}

func (u *Unit) GetNumberOfPoints() int {
	return len(u.Points)
}

func (u *Unit) GetID() int {
	return u.Id
}

func (u *Unit) Dim() int {
	return len(u.Center.Vec)
}

func NewUnit(id int, dimension int, rg Range) Unit{
	unit := Unit{
		Id: id,
		Dimension: dimension,
		neighbour_units: make(map[Range]*Unit),
		Points: make(map[int]Point),
		Range: rg,
	}
	return unit
}