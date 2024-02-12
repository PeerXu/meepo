package meepo_core

import (
	"context"

	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiGetTeleportation(ctx context.Context, req sdk_interface.GetTeleportationRequest) (res sdk_interface.TeleportationView, err error) {
	tp, err := mp.GetTeleportation(ctx, req.TeleportationID)
	if err != nil {
		return
	}

	return ViewTeleportation(tp), nil
}
