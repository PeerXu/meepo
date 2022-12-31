package meepo_core

import (
	"context"

	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPIGetTeleportation(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.GetTeleportationRequest)

	tp, err := mp.GetTeleportation(ctx, req.TeleportationID)
	if err != nil {
		return nil, err
	}

	return ViewTeleportation(tp), nil
}
