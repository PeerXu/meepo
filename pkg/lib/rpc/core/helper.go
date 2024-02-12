package rpc_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type EMPTY = struct{}

var (
	NO_REQUEST EMPTY
	NO_CONTENT EMPTY
)

func WrapHandleFuncGenerics[IT, OT any](fn func(context.Context, IT) (OT, error)) rpc_interface.HandleFunc {
	return func(ctx context.Context, in rpc_interface.HandleRequest) (out rpc_interface.HandleResponse, err error) {
		var req IT
		var res OT

		if err = marshaler.Unmarshal(ctx, in, &req); err != nil {
			return nil, err
		}

		if res, err = fn(ctx, req); err != nil {
			return nil, err
		}

		if out, err = marshaler.Marshal(ctx, res); err != nil {
			return nil, err
		}

		return
	}
}
