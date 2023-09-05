package rpc_http

import (
	"context"

	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func (x *HttpCaller) CallStream(ctx context.Context, method string, opts ...rpc_interface.CallStreamOption) (rpc_interface.Stream, error) {
	panic("unimplemented")
}
