package rpc_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

func empty() any { return &struct{}{} }

var (
	NO_REQUEST = empty
	NO_CONTENT = empty
)

func WrapHandleFunc(newRequest func() any, fn func(context.Context, any) (any, error)) rpc_interface.HandleFunc {
	return func(ctx context.Context, in rpc_interface.HandleRequest) (out rpc_interface.HandleResponse, err error) {
		var res any
		req := newRequest()

		if err = marshaler.Unmarshal(ctx, in, req); err != nil {
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
