package helpers

import "sync"

type hashable interface {
	~string | ~uint32 | ~uint64
}

type Shard[K hashable, V any] struct {
	content map[K]V
	lock sync.RWMutex
}

type SyncMap[K hashable, V any] struct {
	numShards int
	shards []*Shard[K, V]
	hashFn HashFn[K]
}

const DEFAULT_NUM_SHARDS = 32
const DEFAULT_MAP_SIZE = 1024

func NewSyncMap[K hashable, V any](hashFn HashFn[K]) *SyncMap[K, V] {
	shards := make([]*Shard[K, V], DEFAULT_NUM_SHARDS)
	for i := 0; i < DEFAULT_NUM_SHARDS; i++ {
		shards[i] = &Shard[K, V]{
			content: make(map[K]V, DEFAULT_MAP_SIZE / DEFAULT_NUM_SHARDS),
			lock: sync.RWMutex{},
		}
	}

	return &SyncMap[K, V]{
		numShards: DEFAULT_NUM_SHARDS,
		shards: shards,
		hashFn: hashFn,
	}
}

func (m *SyncMap[K, V]) Get(key K) (value V, ok bool) {
	shard := m.hashFn(key) & uint64(m.numShards - 1)

	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].lock.RLock()
	defer m.shards[shard].lock.RUnlock()

	if v, ok := m.shards[shard].content[key]; ok {
		return v, true
	}

	return
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	shard := m.hashFn(key) & uint64(m.numShards - 1)

	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].lock.Lock()
	defer m.shards[shard].lock.Unlock()

	m.shards[shard].content[key] = value
}

func (m *SyncMap[K, V]) Delete(key K) {
	shard := m.hashFn(key) & uint64(m.numShards - 1)

	if m.shards[shard] == nil {
		return
	}

	m.shards[shard].lock.Lock()
	defer m.shards[shard].lock.Unlock()

	if _, ok := m.shards[shard].content[key]; ok {
		delete(m.shards[shard].content, key)
	}
}

func (m *SyncMap[K, V]) Len() int {
	length := 0

	for _, shard := range m.shards {
		shard.lock.RLock()
		length += len(shard.content)
		shard.lock.RUnlock()
	}

	return length
}

func (m *SyncMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())

	for _, shard := range m.shards {
		shard.lock.RLock()
		for key := range shard.content {
			keys = append(keys, key)
		}
		shard.lock.RUnlock()
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
		RWMutex: sync.RWMutex{},
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

type SyncQueue[T any] struct {
	list SyncList[T] // You can replace "interface{}" with the specific type you want to store in the queue
}

func NewSyncQueue[T any]() *SyncQueue[T] {
	return &SyncQueue[T]{
		list: *NewSyncList[T](),
	}
}

func (q *SyncQueue[T]) Enqueue(item T) {
	q.list.Add(item)
}

func (q *SyncQueue[T]) Dequeue() (item T, ok bool) {
	if q.list.Len() > 0 {
		item, ok = q.list.Get(0)
		q.list.Remove(0)
	}
	return
}

func (q *SyncQueue[T]) Len() int {
	return q.list.Len()
}