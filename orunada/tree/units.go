package tree


type Units map[Range]*Unit


func (us Units) AddUnit(rg Range, unit *Unit){
	us[rg] = unit
}

func (us Units) GetNeighbouringUnits(rg Range, interval_l float64) []*Unit {
	/**
	U = unit; n{x} => neighbouring units
	|n3|n5|n8|
	|n2|U |n7|
	|n1|n4|n6|
	 */
	n1 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.Low[1] - interval_l},
		High: [2]float64{rg.Low[0], rg.Low[1]}}
	n2 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.Low[1]},
		High: [2]float64{rg.Low[0], rg.High[1]}}
	n3 := Range{Low:[2]float64{rg.Low[0] - interval_l, rg.High[1]},
		High: [2]float64{rg.Low[0], rg.High[1] + interval_l}}

	n4 := Range{Low:[2]float64{rg.Low[0], rg.Low[1] - interval_l},
		High: [2]float64{rg.High[0], rg.Low[1]}}
	n5 := Range{Low:[2]float64{rg.Low[0], rg.High[1]},
		High: [2]float64{rg.High[0], rg.High[1] + interval_l}}

	n6 := Range{Low:[2]float64{rg.High[0], rg.Low[1] - interval_l},
		High: [2]float64{rg.High[0] + interval_l, rg.Low[1]}}
	n7 := Range{Low:[2]float64{rg.High[0], rg.Low[1]},
		High: [2]float64{rg.High[0] + interval_l, rg.High[1]}}
	n8 := Range{Low:[2]float64{rg.High[0], rg.High[1]},
		High: [2]float64{rg.High[0] + interval_l, rg.High[1] + interval_l}}

	neighbour_units := []*Unit{
		us[n1],
		us[n2],
		us[n3],
		us[n4],
		us[n5],
		us[n6],
		us[n7],
		us[n8],
	}

	return neighbour_units
}