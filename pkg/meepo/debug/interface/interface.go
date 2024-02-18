package meepo_debug_interface

import (
	"context"
)

type MeepoDebugInterface interface {
	TransportStateChange(ctx context.Context, happenedAt, host, target, state string) error
}
