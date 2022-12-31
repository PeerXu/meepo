package rpc_core

import (
	"sync"

	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type Caller = rpc_interface.Caller

type NewCallerFunc func(...NewCallerOption) (Caller, error)

var newCallerFuncs sync.Map

func NewCaller(name string, opts ...NewCallerOption) (Caller, error) {
	v, ok := newCallerFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedCallerFn(name)
	}

	return v.(NewCallerFunc)(opts...)
}

func RegisterNewCallerFunc(name string, fn NewCallerFunc) {
	newCallerFuncs.Store(name, fn)
}
