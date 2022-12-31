package transport_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func WrapHandleFunc(newRequest func() any, fn func(context.Context, any) (any, error)) meepo_interface.HandleFunc {
	return func(ctx context.Context, in meepo_interface.HandleRequest) (out meepo_interface.HandleResponse, err error) {
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

type key string

const (
	transportKey key = "transport"
)

func ContextWithTransport(ctx context.Context, t Transport) context.Context {
	return context.WithValue(ctx, transportKey, t)
}

func ContextGetTransport(ctx context.Context) Transport {
	v := ctx.Value(transportKey)
	if v != nil {
		return v.(Transport)
	}
	return nil
}
