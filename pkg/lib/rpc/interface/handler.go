package rpc_interface

import (
	"context"
)

type HandleRequest = []byte
type HandleResponse = []byte

type HandleFunc = func(context.Context, HandleRequest) (HandleResponse, error)
type HandleStreamFunc = func(context.Context, Stream) error

type Handler interface {
	Handle(method string, fn HandleFunc, opts ...HandleOption)
	HandleStream(method string, fn HandleStreamFunc, opts ...HandleStreamOption)
	Do(ctx context.Context, method string, req HandleRequest) (HandleResponse, error)
	DoStream(ctx context.Context, method string, stm Stream) error
}
