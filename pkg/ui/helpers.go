package ui

func Map[T, U any](ts []T, f func(T) U) []U {
	result := make([]U, len(ts))

	for i := range ts {
		result[i] = f(ts[i])
	}

	return result
}
