package rpc_default

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (h DefaultHandler) Do(ctx context.Context, method string, req rpc_interface.HandleRequest) (rpc_interface.HandleResponse, error) {
	fn, ok := h[method]
	if !ok {
		return nil, rpc_core.ErrUnsupportedHandlerFn(method)
	}

	return fn(ctx, req)
}
