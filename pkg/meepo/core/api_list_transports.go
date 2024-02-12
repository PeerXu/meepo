package meepo_core

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiListTransports(ctx context.Context, _ rpc_core.EMPTY) (res []sdk_interface.TransportView, err error) {
	ts, err := mp.ListTransports(ctx)
	if err != nil {
		return
	}

	return ViewTransports(ts), nil
}
