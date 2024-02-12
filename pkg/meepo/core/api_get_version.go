package meepo_core

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	"github.com/PeerXu/meepo/pkg/lib/version"
)

func (mp *Meepo) apiGetVersion(ctx context.Context, _ rpc_core.EMPTY) (version.V, error) {
	return *version.Get(), nil
}
