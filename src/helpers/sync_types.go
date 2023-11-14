package helpers

import "sync"

type SyncMap[K comparable, V any] struct {
	content map[K]V
	sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		content: make(map[K]V),
		RWMutex: sync.RWMutex{},
	}
}

func (m *SyncMap[K, V]) Get(key K) (value V, ok bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok = m.content[key]
	return
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()
	m.content[key] = value
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.Lock()
	defer m.Unlock()
	delete(m.content, key)
}

func (m *SyncMap[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.content)
}

func (m *SyncMap[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()
	keys := make([]K, 0, len(m.content))
	for key := range m.content {
		keys = append(keys, key)
	}
	return keys
}

type SyncSet[V comparable] struct {
	m map[V]struct{}
	sync.RWMutex
}

func NewSyncSet[V comparable]() *SyncSet[V] {
	return &SyncSet[V]{
		m:       make(map[V]struct{}),
		RWMutex: sync.RWMutex{},
	}
}

func (s *SyncSet[V]) Add(item V) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = struct{}{}
}

func (s *SyncSet[V]) Contains(item V) bool {
	s.RLock()
	defer s.RUnlock()
	_, exists := s.m[item]
	return exists
}

func (s *SyncSet[V]) Remove(item V) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *SyncSet[V]) ToSlice() []V {
	s.RLock()
	defer s.RUnlock()
	items := make([]V, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}

type SyncList[V any] struct {
	content []V
	sync.RWMutex
}

func NewSyncList[V any]() *SyncList[V] {
	return &SyncList[V]{
		content: make([]V, 0),
	}
}

func (l *SyncList[V]) Add(item V) {
	l.Lock()
	defer l.Unlock()
	l.content = append(l.content, item)
}

func (l *SyncList[V]) Get(index int) (item V, ok bool) {
	l.RLock()
	defer l.RUnlock()
	if index >= 0 && index < len(l.content) {
		item = l.content[index]
		ok = true
	}
	return
}

func (l *SyncList[V]) Len() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.content)
}

func (l *SyncList[V]) Remove(index int) {
	l.Lock()
	defer l.Unlock()
	if index >= 0 && index < len(l.content) {
		l.content = append(l.content[:index], l.content[index+1:]...)
	}
}
