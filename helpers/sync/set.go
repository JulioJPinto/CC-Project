package sync

import "sync"

type SyncSet[T any] struct {
	set sync.Map
}

// SyncSet creates a new ParametricSyncSet.
func NewSyncSet[T any]() *SyncSet[T] {
	return &SyncSet[T]{
		set: sync.Map{},
	}
}

// Add adds an element to the set.
func (s *SyncSet[T]) Add(item T) {
	s.set.Store(item, true)
}

// Remove removes an element from the set.
func (s *SyncSet[T]) Remove(item T) {
	s.set.Delete(item)
}

// Contains checks if the set contains a specific element.
func (s *SyncSet[T]) Contains(item T) bool {
	_, ok := s.set.Load(item)
	return ok
}

// Size returns the size of the set.
func (s *SyncSet[T]) Size() int {
	size := 0
	s.set.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

// List returns a slice containing all elements in the set.
func (s *SyncSet[T]) List() []T {
	var list []T
	s.set.Range(func(key, _ interface{}) bool {
		list = append(list, key.(T))
		return true
	})
	return list
}