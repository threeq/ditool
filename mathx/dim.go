package mathx

func Max[V int8 | int | int64 | int32 | float32 | float64](a V, b ...V) V {
	if len(b) == 0 {
		return a
	}
	for i := range b {
		if b[i] > a {
			a = b[i]
		}
	}
	return a
}

func Min[V int8 | int | int64 | int32 | float32 | float64](a V, b ...V) V {
	if len(b) == 0 {
		return a
	}
	for i := range b {
		if b[i] < a {
			a = b[i]
		}
	}
	return a
}
