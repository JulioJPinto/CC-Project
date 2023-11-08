package helpers

import "encoding/json"

type Set[T comparable] struct {
	data []T
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	s.data = make([]T, len(arr))
	for _, item := range arr {
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
	return &Set[T]{make([]T, 0)}
}

func (s *Set[T]) Slice() []T {
	ret := make([]T, 0)
	copy(s.data, ret)
	return ret
}

func (s *Set[T]) Contains(elem T) bool {
	for _, e := range s.data {
		if e == elem {
			return true
		}
	}
	return false
}

func (s *Set[T]) AnyMatch(f func(T) bool) bool {
	for _, item := range s.data {
		if f(item) {
			return true
		}
	}
	return false
}

func (s *Set[T]) Add(elem T) {
	// Check if the element already exists in the set before adding it
	if !s.Contains(elem) {
		s.data = append(s.data, elem)
	}
}

func (s *Set[T]) Remove(elem T){
	for i, e := range s.data {
        if e == elem {
            s.data = append(s.data[:i], s.data[i+1:]...)
        }
    }
}

func (s *Set[T]) RemoveIf(f func(T) bool) {
	for i, e := range s.data {
        if f(e) {
            s.data = append(s.data[:i], s.data[i+1:]...)
        }
    }
}

func (s *Set[T]) Union(other Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, elem := range s.data {
		result.Add(elem)
	}
	for _, elem := range other.data {
		result.Add(elem)
	}
	return result
}

func (s *Set[T]) Intersection(other Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, elem := range s.data {
		if other.Contains(elem) {
			result.Add(elem)
		}
	}
	return result
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
