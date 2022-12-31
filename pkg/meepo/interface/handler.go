package meepo_interface

import "context"

type HandleRequest = []byte
type HandleResponse = []byte

type HandleFunc = func(context.Context, HandleRequest) (HandleResponse, error)

type Handler interface {
	Handle(method string, fn HandleFunc, opts ...HandleOption)
}
