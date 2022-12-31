package tracker_rpc

import (
	"context"
)

func (tk *RPCTracker) context() context.Context {
	return context.Background()
}
