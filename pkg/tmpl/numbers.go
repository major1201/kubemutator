package tmpl

import "math/rand"

func inc(a int) int {
	return a + 1
}

func add(a, b float64) float64 {
	return a + b
}

func sub(a, b float64) float64 {
	return a - b
}

func mul(a, b float64) float64 {
	return a * b
}

func div(a, b float64) float64 {
	return a / b
}

func mod(a, b int64) int64 {
	return a % b
}

func random() int64 {
	return rand.Int63()
}
