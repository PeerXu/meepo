package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPICloseTeleportation(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.CloseTeleportationRequest)

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":         "hdrAPICloseTeleportation",
		"teleportationID": req.TeleportationID,
	})

	tp, err := mp.GetTeleportation(ctx, req.TeleportationID)
	if err != nil {
		logger.WithError(err).Errorf("failed to get teleportation")
		return nil, err
	}

	if err = tp.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close teleportation")
		return nil, err
	}

	logger.Infof("teleportation closed")

	return rpc_core.NO_CONTENT, nil
}
