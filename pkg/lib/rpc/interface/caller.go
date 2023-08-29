package rpc_interface

import "context"

type CallRequest = any

type CallResponse = any

type Caller interface {
	Call(ctx context.Context, method string, req CallRequest, res CallResponse, opts ...CallOption) error
	// CallStream(ctx context.Context, method string, opts ...CallStreamOption) (Stream, error)
}
