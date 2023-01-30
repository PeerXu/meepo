package atomic

import "sync/atomic"

type GenericsValue[T any] interface {
	CompareAndSwap(old, new T) (swapped bool)
	Load() (val T)
	Store(val T)
	Swap(new T) (old T)
}

type genericsValue[T any] struct {
	*atomic.Value
}

func NewValue[T any]() GenericsValue[T] {
	return &genericsValue[T]{&atomic.Value{}}
}

func (v *genericsValue[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.Value.CompareAndSwap(old, new)
}

func (v *genericsValue[T]) Load() (val T) {
	t := v.Value.Load()
	if t == nil {
		return
	}
	return t.(T)
}

func (v *genericsValue[T]) Store(val T) {
	v.Value.Store(val)
}

func (v *genericsValue[T]) Swap(new T) (old T) {
	return v.Value.Swap(new).(T)
}
