package meepo_core

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
)

func (mp *Meepo) apiWhoami(ctx context.Context, _ rpc_core.EMPTY) (string, error) {
	return mp.Addr().String(), nil
}
