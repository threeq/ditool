package mathx

import "math"

func NaN() float64 {
	return math.NaN()
}

func InfMax() float64 {
	return math.Inf(1)
}

func InfMin() float64 {
	return math.Inf(-1)
}
