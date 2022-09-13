package mathx

import "math"

// Max returns the larger of x or y.
//
// Special cases are:
//
//	Max(x, +Inf) = Max(+Inf, x) = +Inf
//	Max(x, NaN) = Max(NaN, x) = NaN
//	Max(+0, ±0) = Max(±0, +0) = +0
//	Max(-0, -0) = -0
func Max[V int8 | int | int32 | int64 | float32 | float64](a V, b ...V) V {
	if len(b) == 0 {
		return a
	}
	for i := range b {
		a = max(a, b[i])
	}
	return a
}

func max[V int8 | int | int64 | int32 | float32 | float64](x, y V) V {
	// special cases
	switch {
	case math.IsInf(float64(x), 1) || math.IsInf(float64(y), 1):
		return V(math.Inf(1))
	case math.IsNaN(float64(x)) || math.IsNaN(float64(y)):
		return V(math.NaN())
	case x == 0 && x == y:
		if math.Signbit(float64(x)) {
			return y
		}
		return x
	}
	if x > y {
		return x
	}
	return y
}

// Min returns the smaller of x or y.
//
// Special cases are:
//
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//	Min(-0, ±0) = Min(±0, -0) = -0
func Min[V int8 | int | int64 | int32 | float32 | float64](a V, b ...V) V {
	if len(b) == 0 {
		return a
	}
	for i := range b {
		a = min(a, b[i])
	}
	return a
}

func min[V int8 | int | int64 | int32 | float32 | float64](x, y V) V {
	// special cases
	switch {
	case math.IsInf(float64(x), -1) || math.IsInf(float64(y), -1):
		return V(math.Inf(-1))
	case math.IsNaN(float64(x)) || math.IsNaN(float64(y)):
		return V(math.NaN())
	case x == 0 && x == y:
		if math.Signbit(float64(x)) {
			return y
		}
		return x
	}
	if x < y {
		return x
	}
	return y
}
