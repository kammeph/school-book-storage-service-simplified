package fp

func Find[TValue any](values []TValue, f func(value TValue) bool) *TValue {
	idx := FindIndex(values, f)
	if idx == -1 {
		return nil
	}
	return &values[idx]
}

func FindIndex[TValue any](values []TValue, f func(value TValue) bool) int {
	for idx, v := range values {
		if f(v) {
			return idx
		}
	}
	return -1
}
