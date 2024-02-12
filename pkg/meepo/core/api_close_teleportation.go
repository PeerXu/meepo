package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiCloseTeleportation(ctx context.Context, req sdk_interface.CloseTeleportationRequest) (res rpc_core.EMPTY, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":         "hdrAPICloseTeleportation",
		"teleportationID": req.TeleportationID,
	})

	tp, err := mp.GetTeleportation(ctx, req.TeleportationID)
	if err != nil {
		logger.WithError(err).Errorf("failed to get teleportation")
		return
	}

	if err = tp.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close teleportation")
		return
	}

	logger.Infof("teleportation closed")

	return rpc_core.NO_CONTENT, nil
}
