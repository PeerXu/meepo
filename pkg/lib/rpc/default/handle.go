package rpc_default

import (
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (h DefaultHandler) Handle(method string, fn rpc_interface.HandleFunc, opts ...rpc_interface.HandleOption) {
	h[method] = fn
}
