package rpc_default

import (
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type DefaultHandler map[string]rpc_interface.HandleFunc

func NewDefaultHandler(opts ...rpc_interface.NewHandlerOption) (rpc_interface.Handler, error) {
	return DefaultHandler(make(map[string]rpc_interface.HandleFunc)), nil
}

func init() {
	rpc_core.RegisterNewHandlerFunc("default", NewDefaultHandler)
}