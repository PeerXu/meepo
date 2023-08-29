package rpc_default

import rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"

func (h DefaultHandler) HandleStream(method string, fn rpc_interface.HandleStreamFunc, opts ...rpc_interface.HandleStreamOption) {
	h[method] = fn
}
