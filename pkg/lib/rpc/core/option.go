package rpc_core

import (
	"github.com/PeerXu/meepo/pkg/internal/option"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

const (
	OPTION_HANDLER = "handler"
	OPTION_CALLER  = "caller"
)

type NewCallerOption = rpc_interface.NewCallerOption

type NewHandlerOption = rpc_interface.NewHandlerOption

type NewServerOption = rpc_interface.NewServerOption

func WithHandler(h Handler) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_HANDLER] = h
	}
}

func GetHandler(o option.Option) (Handler, error) {
	i := o.Get(OPTION_HANDLER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_HANDLER)
	}

	v, ok := i.(Handler)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}

	return v, nil
}

func WithCaller(c Caller) option.ApplyOption {
	return func(o option.Option) {
		o[OPTION_CALLER] = c
	}
}

func GetCaller(o option.Option) (Caller, error) {
	i := o.Get(OPTION_CALLER).Inter()
	if i == nil {
		return nil, option.ErrOptionRequiredFn(OPTION_CALLER)
	}

	v, ok := i.(Caller)
	if !ok {
		return nil, option.ErrUnexpectedTypeFn(v, i)
	}

	return v, nil
}
