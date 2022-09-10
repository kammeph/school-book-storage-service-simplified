package fp

func Some[TValue any](values []TValue, f func(value TValue) bool) bool {
	for _, v := range values {
		if f(v) {
			return true
		}
	}
	return false
}
