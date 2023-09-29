package atomic

import "sync/atomic"

type GenericValue[T any] interface {
	CompareAndSwap(old, new T) (swapped bool)
	Load() (val T)
	Store(val T)
	Swap(new T) (old T)
}

type genericValue[T any] struct {
	*atomic.Value
}

func NewValue[T any]() GenericValue[T] {
	return &genericValue[T]{&atomic.Value{}}
}

func (v *genericValue[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.Value.CompareAndSwap(any(old), any(new))
}

func (v *genericValue[T]) Load() (val T) {
	t := v.Value.Load()
	if t == nil {
		return
	}
	return t.(T)
}

func (v *genericValue[T]) Store(val T) {
	v.Value.Store(val)
}

func (v *genericValue[T]) Swap(new T) (old T) {
	return v.Value.Swap(new).(T)
}
