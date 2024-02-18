package lib_registerer

import "github.com/PeerXu/meepo/pkg/lib/option"

type RegisterFunc[T any] func(name string, fn NewInstanceFunc[T])
type NewFunc[T any] func(name string, opts ...option.ApplyOption) (T, error)

func Pair[T any]() (RegisterFunc[T], NewFunc[T]) {
	r := NewRegisterer[T]()
	return r.Register, r.New
}
