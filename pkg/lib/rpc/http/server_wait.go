package rpc_http

import "context"

func (s *HttpServer) Wait(context.Context) <-chan error {
	return s.errors
}
