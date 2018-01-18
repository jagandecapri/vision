package utils

import "math"

func Gaussian (height float64, x float64, center float64, width float64) float64{
	return height * math.Exp(-math.Pow(x - center, 2)/(2.0*width*width))
}