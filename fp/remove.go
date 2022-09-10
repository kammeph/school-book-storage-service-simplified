package fp

func Remove[TValue any](values []TValue, f func(value TValue) bool) []TValue {
	idx := FindIndex(values, f)
	return RemoveIndex(values, idx)
}

func RemoveIndex[TValue any](values []TValue, idx int) []TValue {
	if idx == -1 {
		return values
	}
	return append(values[:idx], values[idx+1:]...)
}

// func RemoveRef[TValue any](values *[]TValue, f func(value TValue) bool) []TValue {
// 	for idx, v := range *values {
// 		if f(v) {
// 			return append(values[:idx], values[idx+1:]...)
// 		}
// 	}
// 	return values
// }
