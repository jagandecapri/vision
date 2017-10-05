package tree

import "testing"

func TestCluster(t *testing.T) {
	units := Units{Store: make(map[Range]*Unit)}
	interval_l := 1.0

	u1 := Unit{Id: 1}
	r1 := Range{Low: [2]float64{0, 0}, High: [2]float64{1, 1}}
	units.AddUnit(&u1, r1, interval_l)

	u2 := Unit{Id: 2}
	r2 := Range{Low: [2]float64{0, 1}, High: [2]float64{1, 2}}
	units.AddUnit(&u2, r2, interval_l)

	u3 := Unit{Id: 3}
	r3 := Range{Low: [2]float64{0, 2}, High: [2]float64{1, 3}}
	units.AddUnit(&u3, r3, interval_l)

	u4 := Unit{Id: 4}
	r4 := Range{Low: [2]float64{1, 0}, High: [2]float64{2, 1}}
	units.AddUnit(&u4, r4, interval_l)

	u5 := Unit{Id: 5}
	r5 := Range{Low: [2]float64{1, 2}, High: [2]float64{2, 3}}
	units.AddUnit(&u5, r5, interval_l)

	u6 := Unit{Id: 6}
	r6 := Range{Low: [2]float64{2, 0}, High: [2]float64{3, 1}}
	units.AddUnit(&u6, r6, interval_l)

	u7 := Unit{Id: 7}
	r7 := Range{Low: [2]float64{2, 1}, High: [2]float64{3, 2}}
	units.AddUnit(&u7, r7, interval_l)

	u8 := Unit{Id: 8}
	r8 := Range{Low: [2]float64{2, 2}, High: [2]float64{3, 3}}
	units.AddUnit(&u8, r8, interval_l)

	u9 := Unit{Id: 9}
	r9 := Range{Low: [2]float64{1, 1}, High: [2]float64{2, 2}}
	units.AddUnit(&u9, r9, interval_l)

}
