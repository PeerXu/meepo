package meepo_core

import (
	"context"

	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiListTeleportations(ctx context.Context, _ rpc_core.EMPTY) (res []sdk_interface.TeleportationView, err error) {
	tps, err := mp.ListTeleportations(ctx)
	if err != nil {
		return nil, err
	}

	return ViewTeleportations(tps), nil
}
