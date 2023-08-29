package rpc_core

import "github.com/PeerXu/meepo/pkg/lib/errors"

var (
	ErrUnsupportedCaller, ErrUnsupportedCallerFn   = errors.NewErrorAndErrorFunc[string]("unsupported caller")
	ErrUnsupportedServer, ErrUnsupportedServerFn   = errors.NewErrorAndErrorFunc[string]("unsupported server")
	ErrUnsupportedHandler, ErrUnsupportedHandlerFn = errors.NewErrorAndErrorFunc[string]("unsupported handler")
	ErrUnsupportedMethod, ErrUnsupportedMethodFn   = errors.NewErrorAndErrorFunc[string]("unsupported method")
)
