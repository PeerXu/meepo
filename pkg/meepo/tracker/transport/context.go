package tracker_transport

import (
	"context"
)

func (tk *TransportTracker) context() context.Context {
	return context.Background()
}
