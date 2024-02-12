package rpc_core

import "context"

type Middleware[IT, OT any] func(next func(context.Context, IT) (OT, error)) func(context.Context, IT) (OT, error)
