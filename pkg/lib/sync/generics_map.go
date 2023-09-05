package sync

import "sync"

type GenericMap[K any, V any] interface {
	Store(key K, value V)
	Delete(key K)
	Load(key K) (value V, ok bool)
	LoadAndDelete(key K) (value V, loaded bool)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Range(f func(key K, value V) bool)
}

type genericMap[K, V any] struct {
	*sync.Map
}

func NewMap[K any, V any]() GenericMap[K, V] {
	return &genericMap[K, V]{&sync.Map{}}
}

func (m *genericMap[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

func (m *genericMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.Map.Load(key)
	if !ok {
		return
	}
	return v.(V), ok
}

func (m *genericMap[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

func (m *genericMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.Map.LoadAndDelete(key)
	if !loaded {
		return
	}
	return v.(V), loaded
}

func (m *genericMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	if !loaded {
		return
	}
	return a.(V), loaded
}

func (m *genericMap[K, V]) Range(f func(K, V) bool) {
	m.Map.Range(func(key any, value any) bool {
		return f(key.(K), value.(V))
	})
}
