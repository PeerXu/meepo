package rpc_simple_http

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/marshaler"
)

func (c *SimpleHttpCaller) context() context.Context {
	return marshaler.ContextWithMarshalerAndUnmarshaler(context.Background(), c.marshaler, c.unmarshaler)
}

func (s *SimpleHttpServer) context() context.Context {
	return marshaler.ContextWithMarshalerAndUnmarshaler(context.Background(), s.marshaler, s.unmarshaler)
}
