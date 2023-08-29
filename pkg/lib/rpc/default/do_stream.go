package rpc_default

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (h DefaultHandler) DoStream(ctx context.Context, method string, stm rpc_interface.Stream) error {
	fn, ok := h[method]
	if !ok {
		return rpc_core.ErrUnsupportedMethodFn(method)
	}

	return fn.(rpc_interface.HandleStreamFunc)(ctx, stm)
}
