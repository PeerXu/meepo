package transport_core

import (
	"sync"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

type Transport = meepo_interface.Transport

type NewTransportFunc func(...NewTransportOption) (Transport, error)

var newTransportFuncs sync.Map

func NewTransport(name string, opts ...NewTransportOption) (Transport, error) {
	v, ok := newTransportFuncs.Load(name)
	if !ok {
		return nil, ErrUnsupportedTransportFn(name)
	}

	return v.(NewTransportFunc)(opts...)
}

func RegisterNewTransportFunc(name string, fn NewTransportFunc) {
	newTransportFuncs.Store(name, fn)
}
