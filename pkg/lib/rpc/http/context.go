package rpc_http

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/marshaler"
)

func (x *HttpServer) context() context.Context {
	return marshaler.ContextWithMarshalerAndUnmarshaler(context.Background(), x.marshaler, x.unmarshaler)
}
