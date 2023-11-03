package helpers

type Set[T any] struct {
	slice []T
	eq    func(a T, b T) bool
}

func NewSet[T any](func(a T, b T) bool) *Set[T] {
	return &Set[T]{}
}

func NewSetFromSlice[T comparable](v []T) *Set[T] {
	ret := NewDefaultSet[T]()
	for _, item := range v {
		ret.Add(item)
	}
	return ret
}

func NewDefaultSet[T comparable]() *Set[T] {
	eq := func(a T, b T) bool { return a == b }
	return &Set[T]{make([]T, 0), eq}
}

func (s *Set[T]) Slice() []T {
	ret := make([]T, 0)
	copy(s.slice, ret)
	return ret
}

func (s *Set[T]) Contains(elem T) bool {
	for _, e := range s.slice {
		if s.eq(e, elem) {
			return true
		}
	}
	return false
}

func (s *Set[T]) AnyMatch(f func(T) bool) bool {
	for _, item := range s.slice {
		if f(item) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(elem T) {
	// Check if the element already exists in the set before adding it
	if !s.Contains(elem) {
		s.slice = append(s.slice, elem)
	}
}

func (s *Set[T]) Union(other Set[T]) *Set[T] {
	result := NewSet[T](s.eq)
	for _, elem := range s.slice {
		result.Add(elem)
	}
	for _, elem := range other.slice {
		result.Add(elem)
	}
	return result
}

func (s *Set[T]) Intersection(other Set[T]) *Set[T] {
	result := NewSet[T](s.eq)
	for _, elem := range s.slice {
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