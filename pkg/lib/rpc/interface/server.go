package rpc_interface

import "context"

type Server interface {
	Serve(context.Context) <-chan error
	Terminate(context.Context) error
	Wait(context.Context) <-chan error
}
