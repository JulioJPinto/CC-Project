package helpers

func MapKeys[T comparable, T2 any](m map[T]T2) *Set[T] {
	keys := make([]T, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return NewSetFromSlice(keys)
}

func MapVals[T any](m map[any]T) []T {
	vals := make([]T, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

/*
m1 is altered, gains the keys and values of m2

In case of colision, the value of m2 is used
*/
func MergeMaps[T comparable, U any](m1 map[T]U, m2 map[T]U) {
	for k, v := range m2 {
		m1[k] = v
	}
}
