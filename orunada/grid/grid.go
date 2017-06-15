package grid

import (
	"github.com/cockroachdb/apd"
	"fmt"
)

type Grid struct{
	units []Unit
	axes [2]string
}

func (g *Grid) Build2DGrid(axes []string, ctx *apd.Context) *Grid{
	copy(g.axes[:], axes)
	//ctx := apd.BaseContext.WithPrecision(6)
	dim := []Interval{}
	for range axes {
		interval_l, _ := new(apd.Decimal).SetFloat64(0.1)
		lb, _ := new(apd.Decimal).SetFloat64(0.0)
		ub, _ := new(apd.Decimal).SetFloat64(0.1)
		for i := 0; i < 10; i++ {
			interval := Interval{}
			lb_tmp, _ := lb.Float64()
			ub_tmp, _ := ub.Float64()
			interval.range_ = []float64{lb_tmp, ub_tmp}
			dim = append(dim, interval)
			ctx.Add(lb, lb, interval_l)
			ctx.Add(ub, ub, interval_l)
		}
	}

	unit_id := 0
	for i := 0; i < len(dim); i++{
		for j := 0; j < len(dim); j++{
			unit := Unit{}
			unit_id += 1
			unit.id = unit_id
			unit.intervals = append(unit.intervals, dim[i], dim[j])
			g.units = append(g.units, unit)
		}
	}
	return g
}

func (g *Grid) intersect(p Point) *Unit{
	vec := p.norm_vec
	ctx := apd.BaseContext.WithPrecision(6)
	for i := 0; i < len(g.units); i++{
		unit := &g.units[i]
		inside_interval_ctr := false
		lower_bound_x := unit.intervals[0].range_[0]
		upper_bound_x := unit.intervals[0].range_[1]
		lower_bound_y := unit.intervals[1].range_[0]
		upper_bound_y := unit.intervals[1].range_[1]

		lb_x, _ := new(apd.Decimal).SetFloat64(lower_bound_x)
		ub_x, _ := new(apd.Decimal).SetFloat64(upper_bound_x)
		lb_y, _ := new(apd.Decimal).SetFloat64(lower_bound_y)
		ub_y, _ := new(apd.Decimal).SetFloat64(upper_bound_y)

		vec_0, _ := new(apd.Decimal).SetFloat64(vec[g.axes[0]])
		vec_1, _ := new(apd.Decimal).SetFloat64(vec[g.axes[1]])

		cmp_vec_0_lb_x := new(apd.Decimal)
		cmp_vec_0_ub_x := new(apd.Decimal)
		ctx.Cmp(cmp_vec_0_lb_x, vec_0, lb_x)
		ctx.Cmp(cmp_vec_0_ub_x, vec_0, ub_x)

		cmp_vec_1_lb_y := new(apd.Decimal)
		cmp_vec_1_ub_y := new(apd.Decimal)
		ctx.Cmp(cmp_vec_1_lb_y, vec_1, lb_y)
		ctx.Cmp(cmp_vec_1_ub_y, vec_1, ub_y)

		int_cmp_vec_0_lb_x, _ := cmp_vec_0_lb_x.Int64()
		int_cmp_vec_0_ub_x, _ := cmp_vec_0_ub_x.Int64()
		int_cmp_vec_1_lb_y, _ := cmp_vec_1_lb_y.Int64()
		int_cmp_vec_1_ub_y, _ := cmp_vec_1_ub_y.Int64()

		if i == len(g.units) - 1{
			if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1 || int_cmp_vec_0_ub_x == 0) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1 || int_cmp_vec_1_ub_y == 0){
				inside_interval_ctr = true
			}
		} else {
			if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1){
				inside_interval_ctr = true
			}
		}

		if inside_interval_ctr == true{
			fmt.Println("Intersected", unit)
			return unit
		}
	}
	return nil
}

func (g *Grid) calculateListOfDenseUnits(){

}