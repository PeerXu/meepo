package rpc_core

import (
	"sync"

	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type Handler = rpc_interface.Handler

type NewHandlerFunc func(...NewHandlerOption) (Handler, error)

var newHandlerFuncs sync.Map

func NewHandler(name string, opts ...NewHandlerOption) (Handler, error) {
	v, ok := newHandlerFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedHandlerFn(name)
	}

	return v.(NewHandlerFunc)(opts...)
}

func RegisterNewHandlerFunc(name string, fn NewHandlerFunc) {
	newHandlerFuncs.Store(name, fn)
}
