package sync

import "sync"

type Map[K comparable,V any] struct {
	internalMap sync.Map
}
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	result, ok := m.internalMap.Load(key)
	if ok {
		value = result.(V)
	}
	return value, ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.internalMap.Store(key, value)
}

func (m *Map[K, V]) Delete(key K) {
	m.internalMap.Delete(key)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.internalMap.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) Fold(initialValue interface{}, folder func(accumulator interface{}, key K, value V) interface{}) interface{} {
	accumulator := initialValue

	m.Range(func(key K, value V) bool {
		accumulator = folder(accumulator, key, value)
		return true
	})

	return accumulator
}
