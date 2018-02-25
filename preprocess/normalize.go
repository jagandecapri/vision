package preprocess

import (
	"github.com/jagandecapri/vision/tree"
)

type DimMinMax struct{
	Min, Max float64
	Range float64
}

func norm_mat(elem float64, col_min float64, col_max float64) float64{
	return (elem - col_min)/(col_max - col_min)
}

func Normalize(mat []tree.Point, sorter []string) []tree.Point {
	rows := len(mat)

	dim_min_max := map[string]DimMinMax{}
	for _,c := range sorter {
		min := mat[0].Vec_map[c]
		max := mat[0].Vec_map[c]
		for j := 0; j < rows; j++{
			val := mat[j].Vec_map[c]
			if val < min{
				min = val
			} else if  val > max{
				max = val
			}
		}
		range_ := max - min
		dim_min_max[c] = DimMinMax{min, max, range_}
	}

	for i := 0; i < rows; i++{
		for _, c := range sorter{
			col_min := dim_min_max[c].Min
			col_max := dim_min_max[c].Max
			elem := mat[i].Vec_map[c]
			if col_min == 0 && col_max == 0{
				mat[i].Vec_map[c] = elem
			} else {
				tmp := norm_mat(elem, col_min, col_max)
				if tmp == 1.0{
					mat[i].Vec_map[c] = tmp - 0.0000001 //Trick for points having highs at 1.0
				} else {
					mat[i].Vec_map[c] = norm_mat(elem, col_min, col_max)
				}
			}
		}
	}

	return mat
}
