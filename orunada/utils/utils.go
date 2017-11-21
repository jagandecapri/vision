package utils

import (
	"time"
	"log"
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
func Round(x float64, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
