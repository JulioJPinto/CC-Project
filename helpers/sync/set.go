package sync

type Set[T comparable] struct {
	set Map[T, any]
}

// SyncSet creates a new ParametricSyncSet.
func NewSyncSet[T comparable]() *Set[T] {
	return &Set[T]{
		set: Map[T, any]{},
	}
}

// Add adds an element to the set.
func (s *Set[T]) Add(item T) {
	s.set.Store(item, true)
}

// Remove removes an element from the set.
func (s *Set[T]) Remove(item T) {
	s.set.Delete(item)
}

// Contains checks if the set contains a specific element.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.set.Load(item)
	return ok
}

// Size returns the size of the set.
func (s *Set[T]) Size() int {
	size := 0
	s.set.Range(func(_ T, _ any) bool {
		size++
		return true
	})
	return size
}

func (s *Set[T]) AnyMatch(predicate func(T) bool) bool {
	result := false
	s.set.Fold(false, func(_ any, element T, _ any) interface{} {
		if predicate(element) {
			result = true
			return true // Stop folding since we found a match
		}
		return false
	})
	return result

}

func (s *Set[T]) RemoveIf(predicate func(T) bool) {
	s.set.Range(func(element T, _ any) bool {
		if predicate(element) {
			s.Remove(element)
		}
		return true
	})
}

// List returns a slice containing all elements in the set.
func (s *Set[T]) List() []T {
	var list []T
	s.set.Range(func(key T, _ any) bool {
		list = append(list, key)
		return true
	})
	return list
}
