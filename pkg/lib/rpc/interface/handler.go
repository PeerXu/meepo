package rpc_interface

import (
	"context"
)

type HandleRequest = []byte
type HandleResponse = []byte

type HandleFunc = func(context.Context, HandleRequest) (HandleResponse, error)

type Handler interface {
	Handle(method string, fn HandleFunc, opts ...HandleOption)
	Do(ctx context.Context, method string, req HandleRequest) (HandleResponse, error)
}
