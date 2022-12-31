package rpc_core

import (
	"sync"

	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type Server = rpc_interface.Server

type NewServerFunc func(...NewServerOption) (Server, error)

var newServerFuncs sync.Map

func NewServer(name string, opts ...NewServerOption) (Server, error) {
	v, ok := newServerFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedServerFn(name)
	}

	return v.(NewServerFunc)(opts...)
}

func RegisterNewServerFunc(name string, fn NewServerFunc) {
	newServerFuncs.Store(name, fn)
}
