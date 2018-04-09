package utils

import (
	"time"
	"go.uber.org/zap"
	"math"
)

func Comb(n, m int, emit func([]int)) {
	s := make([]int, m)
	last := m - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				emit(s)
			} else {
				rc(i+1, j+1)
			}
		}
		return
	}
	rc(0, 0)
}

func UniqFloat64(input []float64) []float64 {
	u := make([]float64, 0, len(input))
	m := make(map[float64]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func UniqInt(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func UniqString(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

func GetKeyComb(sorter []string, feature_cnt int) [][]string {
	all := [][]string{}
	Comb(len(sorter), feature_cnt, func (c []int){
		tmp := []string{}
		for _, v := range c {
			tmp = append(tmp, sorter[v])
		}
		all = append(all, tmp)
	})
	return all
}

//Taken from https://stackoverflow.com/questions/39544571/golang-round-to-nearest-0-05
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func TimeTrack(start time.Time, name string, num_CPU int, logger *zap.Logger) {
	elapsed := time.Since(start)
	logger.Info("Log",
		zap.String("method_name", name),
		zap.Int("num_cpu", num_CPU),
		zap.Duration("elapsed_time", elapsed),
	)
}

func CompareStringSlice(a, b []string) []string {
	for i := len(a) - 1; i >= 0; i-- {
		for _, vD := range b {
			if a[i] == vD {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}
	return a
}
