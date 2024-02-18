package lib_registerer

import (
	"sync"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

type NewInstanceFunc[T any] func(...option.ApplyOption) (T, error)

type Registerer[T any] interface {
	Register(name string, fn NewInstanceFunc[T])
	New(name string, opts ...option.ApplyOption) (T, error)
}

func NewRegisterer[T any]() Registerer[T] {
	return &registerer[T]{}
}

type registerer[T any] struct {
	sync.Map
}

func (r *registerer[T]) Register(name string, fn NewInstanceFunc[T]) {
	r.Store(name, fn)
}

func (r *registerer[T]) New(name string, opts ...option.ApplyOption) (T, error) {
	var inst T

	v, ok := r.Load(name)
	if !ok {
		return inst, ErrUnsupportedFn(name)
	}

	return v.(NewInstanceFunc[T])(opts...)
}
