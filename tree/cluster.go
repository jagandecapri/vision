package tree

import (
	"log"
)

type Cluster struct{
	Cluster_id int
	Num_of_points int
	ListOfUnits map[Range]*Unit
}

func (c *Cluster) GetCenter() Point {
	var Center_vec []float64

	for _, unit := range c.ListOfUnits{
		if len(unit.Points) > 0 && Center_vec == nil{
			for _, p := range unit.Points{
				Center_vec = make([]float64, p.Dim())
				break
			}
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

func ValidateCluster(c Cluster) bool{
	ret := true
	if len(c.ListOfUnits) == 1{
		ret = true
	} else{
		for _, unit := range c.ListOfUnits{
			neighbour_present := false
			for _, neigh_unit := range unit.GetNeighbouringUnits(){
				_, ok := c.ListOfUnits[neigh_unit.Range]
				if ok{
					neighbour_present = true
					break
				}
			}
			if neighbour_present == false{
				log.Printf("Neighbour not present for cluster: %v \n", c.Cluster_id)
				for _, unit := range c.ListOfUnits{
					log.Printf("Unit Id: %v Range: %+v", unit.Id, unit.Range)
					for _, neigh_unit := range unit.GetNeighbouringUnits(){
						log.Printf("\tUnit Id: %v Range: %+v", neigh_unit.Id, neigh_unit.Range)
					}
					log.Printf("\n")
				}
				ret = false
				break
			}
		}
	}
	return ret
}