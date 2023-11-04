package helpers

type Set[T comparable] struct {
	SetData []T `json:"SetData"`
}

func NewSetFromSlice[T comparable](v []T) *Set[T] {
	ret := NewSet[T]()
	for _, item := range v {
		ret.Add(item)
	}
	return ret
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{make([]T, 0)}
}

func (s *Set[T]) Slice() []T {
	ret := make([]T, 0)
	copy(s.SetData, ret)
	return ret
}

func (s *Set[T]) Contains(elem T) bool {
	for _, e := range s.SetData {
		if e == elem {
			return true
		}
	}
	return false
}

func (s *Set[T]) AnyMatch(f func(T) bool) bool {
	for _, item := range s.SetData {
		if f(item) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(elem T) {
	// Check if the element already exists in the set before adding it
	if !s.Contains(elem) {
		s.SetData = append(s.SetData, elem)
	}
}

func (s *Set[T]) Union(other Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, elem := range s.SetData {
		result.Add(elem)
	}
	for _, elem := range other.SetData {
		result.Add(elem)
	}
	return result
}

func (s *Set[T]) Intersection(other Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, elem := range s.SetData {
		if other.Contains(elem) {
			result.Add(elem)
		}
	}
	return result
}

func MapKeys[T comparable](m map[T]any) *Set[T] {
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

func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
