package rpc_simple_http

import "context"

func (s *SimpleHttpServer) Wait(context.Context) <-chan error {
	return s.errors
}
