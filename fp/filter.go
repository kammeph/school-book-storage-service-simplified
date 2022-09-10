package fp

func Filter[TValue any](values []TValue, f func(value TValue) bool) []TValue {
	result := []TValue{}
	for _, v := range values {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
