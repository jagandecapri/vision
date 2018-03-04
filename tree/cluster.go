package tree

type Cluster struct{
	Cluster_id int
	Num_of_points int
	ListOfUnits map[Range]*Unit
}

func (c *Cluster) GetCenter() Point {
	var Center_vec []float64

	for _, unit := range c.ListOfUnits{
		if Center_vec == nil{
			Center_vec = make([]float64, unit.Dim())
		}
		for _, p := range unit.Points {
			for i := 0; i < p.Dim(); i++ {
				Center_vec[i] = Center_vec[i] + p.GetValue(i)
			}
		}
	}

	for i, _ := range Center_vec {
		Center_vec[i] = Center_vec[i] / float64(c.Num_of_points)
	}

	pc := Point{Vec: Center_vec}
	return pc
}

func (c *Cluster) GetUnits()map[Range]*Unit{
	return c.ListOfUnits
}

func (c *Cluster) GetNumberOfPoints() int{
	tmp := 0
	for _, unit := range c.ListOfUnits{
		tmp += unit.GetNumberOfPoints()
	}
	return tmp
}