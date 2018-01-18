package utils

import "math"

func Gaussian (height float64, x float64, center float64, width float64) float64{
	return height * math.Exp(-math.Pow(x - center, 2)/(2.0*width*width))
}

// Smooth the data using a gaussian kernel.
func GaussianSmooth(data []float64, n float64) []float64{
	smoothed := []float64{};

	for i := 0; i < len(data); i++ {
		i_float64 := float64(i)
		startIdx := math.Max(0, i_float64 - n);
		endIdx := math.Min(float64(len(data)) - 1, i_float64 + n);

		sumWeights := 0.0;
		sumIndexWeight := 0.0;

			for j := startIdx; j < endIdx + 1; j++ {
				indexScore := math.Abs(j - i_float64)/n;
				indexWeight := Gaussian(indexScore, 1, 0, 1);
				sumWeights += (indexWeight * data[int(j)]);
				sumIndexWeight += indexWeight;
			}
			smoothed = append(smoothed, sumWeights/sumIndexWeight)
	}
	return smoothed;
}