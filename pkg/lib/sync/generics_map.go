package sync

import "sync"

type GenericsMap[T any] interface {
	Store(key any, value T)
	Delete(key any)
	Load(key any) (value T, ok bool)
	LoadAndDelete(key any) (value T, loaded bool)
	LoadOrStore(key any, value T) (actual T, loaded bool)
	Range(f func(key any, value T) bool)
}

type genericsMap[T any] struct {
	*sync.Map
}

func NewMap[T any]() GenericsMap[T] {
	return &genericsMap[T]{&sync.Map{}}
}

func (m *genericsMap[T]) Store(key any, value T) {
	m.Map.Store(key, value)
}

func (m *genericsMap[T]) Load(key any) (value T, ok bool) {
	v, ok := m.Map.Load(key)
	if !ok {
		return
	}
	return v.(T), ok
}

func (m *genericsMap[T]) LoadAndDelete(key any) (value T, loaded bool) {
	v, loaded := m.Map.LoadAndDelete(key)
	if !loaded {
		return
	}
	return v.(T), loaded
}

func (m *genericsMap[T]) LoadOrStore(key any, value T) (actual T, loaded bool) {
	a, loaded := m.Map.LoadOrStore(key, value)
	if !loaded {
		return
	}
	return a.(T), loaded
}

func (m *genericsMap[T]) Range(f func(any, T) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key, value.(T))
	})
}
