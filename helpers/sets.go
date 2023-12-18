package helpers

import "encoding/json"

type Set[T comparable] struct {
	data map[T]struct{}
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	var slice []T
	for elem := range s.data {
		slice = append(slice, elem)
	}
	return json.Marshal(slice)
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var slice []T
	if err := json.Unmarshal(data, &slice); err != nil {
		return err
	}

	s.data = make(map[T]struct{})
	for _, item := range slice {
		s.Add(item)
	}

	return nil
}

func NewSetFromSlice[T comparable](v []T) *Set[T] {
	ret := NewSet[T]()
	for _, item := range v {
		ret.Add(item)
	}
	return ret
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{make(map[T]struct{})}
}

func (s *Set[T]) Slice() []T {
	ret := make([]T, 0, len(s.data))
	for elem := range s.data {
		ret = append(ret, elem)
	}
	return ret
}

func (s *Set[T]) Contains(elem T) bool {
	_, exists := s.data[elem]
	return exists
}

func (s *Set[T]) AnyMatch(f func(T) bool) bool {
	for elem := range s.data {
		if f(elem) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(elem T) {
	s.data[elem] = struct{}{}
}

func (s *Set[T]) Remove(elem T) {
	delete(s.data, elem)
}

func (s *Set[T]) RemoveIf(f func(T) bool) {
	for elem := range s.data {
		if f(elem) {
			delete(s.data, elem)
		}
	}
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for elem := range s.data {
		result.Add(elem)
	}
	for elem := range other.data {
		result.Add(elem)
	}
	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for elem := range s.data {
		if other.Contains(elem) {
			result.Add(elem)
		}
	}
	return result
}

func (s *Set[T]) IsSubset(other *Set[T]) bool {
	for elem := range s.data {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
