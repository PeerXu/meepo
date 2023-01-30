package sync

import "sync"

type GenericsMap[K any, V any] interface {
	Store(key K, value V)
	Delete(key K)
	Load(key K) (value V, ok bool)
	LoadAndDelete(key K) (value V, loaded bool)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Range(f func(key K, value V) bool)
}

type genericsMap[K, V any] struct {
	*sync.Map
}

func NewMap[K any, V any]() GenericsMap[K, V] {
	return &genericsMap[K, V]{&sync.Map{}}
}

func (m *genericsMap[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

func (m *genericsMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.Map.Load(key)
	if !ok {
		return
	}
	return v.(V), ok
}

func (m *genericsMap[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

func (m *genericsMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.Map.LoadAndDelete(key)
	if !loaded {
		return
	}
	return v.(V), loaded
}

func (m *genericsMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	if !loaded {
		return
	}
	return a.(V), loaded
}

func (m *genericsMap[K, V]) Range(f func(K, V) bool) {
	m.Map.Range(func(key any, value any) bool {
		return f(key.(K), value.(V))
	})
}
