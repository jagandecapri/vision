package grid

import (
	"github.com/cockroachdb/apd"
	"fmt"
)

type Grid struct{
	units []Unit
	axes [2]string
	dim_min float64
	dim_max float64
}

func (g *Grid) Build2DGrid(axes []string, ctx *apd.Context) *Grid{
	copy(g.axes[:], axes)
	g.dim_min = 0.0
	g.dim_max = 1.0
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
	vec := p.Norm_vec
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

		tmp1 := vec[g.axes[0]]
		tmp2 := vec[g.axes[1]]
		vec_0, _ := new(apd.Decimal).SetFloat64(tmp1)
		vec_1, _ := new(apd.Decimal).SetFloat64(tmp2)

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

		if tmp1 == 1.0 || tmp2 == 1.0{
			if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1 || int_cmp_vec_0_ub_x == 0) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1 || int_cmp_vec_1_ub_y == 0){
				inside_interval_ctr = true
			}
		} else {
			if (int_cmp_vec_0_lb_x == 1 || int_cmp_vec_0_lb_x == 0) && (int_cmp_vec_0_ub_x == -1) && (int_cmp_vec_1_lb_y == 1 || int_cmp_vec_1_lb_y == 0) && (int_cmp_vec_1_ub_y == -1){
				inside_interval_ctr = true
			}
		}

		if inside_interval_ctr == true{
			return unit
		}
	}
	return nil
}

/* Return updated points */
func (g *Grid) Assign(pts []Point){
	for _,p := range pts{
		u := g.intersect(p)
		//This will fail is u happens to be nil, no intersection which should NEVER be the case
		if u != nil{
			p.Unit_id = u.id
			u.points = append(u.points, p)
		} else {
			fmt.Println("INTERSECTION FAILED",p)
		}
	}
}

func (g *Grid) calculateListOfDenseUnits(){

}