package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	matrix              map[int]map[int]bool
	errEmpty            = errors.New("Matrix is empty")
	errCouldNotConverge = errors.New("Could not converge")
	MAX_TRIES           = 20
	QUALITY_NEEDED      = 0.0001
)

func maxValue(mat map[int]map[int]bool) int {
	max := 0
	for v := range mat {
		if v > max {
			max = v
		}
	}
	return max + 1
}

func multiply(mat map[int]map[int]bool, vec []float64) []float64 {

	res := make([]float64, len(vec))

	for out := 0; out < len(vec); out += 1 {

		for j, v := range vec {
			if mat[j][out] {
				res[out] += v
			}
		}
	}

	return res
}

func normalize(vec []float64) {
	var max float64

	for _, v := range vec {
		max = math.Max(max, v)
	}

	for i := range vec {
		vec[i] = vec[i] / max
	}
}

func compare(vec1 []float64, vec2 []float64) (float64, float64) {
	res := make([]float64, len(vec1))

	for i, v := range vec1 {
		res[i] = vec2[i] / v
	}

	var min, max, avg float64
	min = math.MaxFloat64

	for _, v := range res {
		min = math.Min(v, min)
		max = math.Max(v, max)
		avg += v
	}

	avg = avg / float64(len(res))

	return avg, max - min
}

func GetCentrality(mat map[int]map[int]bool, i int) (float64, error) {

	size := maxValue(mat)
	if size == 0 {
		return 0, errEmpty
	}

	vec := make([]float64, size)
	for i := range vec {
		vec[i] = 0
	}
	vec[0] = 1

	for tries := MAX_TRIES; tries > 0; tries -= 1 {
		next := multiply(mat, vec)
		normalize(next)

		mid, dist := compare(next, multiply(mat, next))
		if dist < QUALITY_NEEDED {
			fmt.Println("STEP")
			var res float64
			for target, _ := range mat[i] {
				res += next[target]
			}
			return res / mid, nil
		}
		vec = next
	}

	return 0, errCouldNotConverge
}

func set(mat map[int]map[int]bool, i, j int) {
	if mat[i] == nil {
		mat[i] = make(map[int]bool)
	}
	mat[i][j] = true
	if mat[j] == nil {
		mat[j] = make(map[int]bool)
	}
	mat[j][i] = true
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	lineNo := 0

	mat := make(map[int]map[int]bool)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\n")
		for x, st := range strings.Split(line, " ") {
			i, err := strconv.Atoi(st)
			if err != nil {
				fmt.Println("Format error")
			}
			if i > 0 {
				set(mat, lineNo, x)
			}
		}
		lineNo += 1

	}

	fmt.Println(mat)

	for i := 0; i < lineNo; i += 1 {
		a, b := GetCentrality(mat, i)
		fmt.Println("RES", a, b)
	}

}
