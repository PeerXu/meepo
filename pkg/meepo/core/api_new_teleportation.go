package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiNewTeleportation(ctx context.Context, req sdk_interface.NewTeleportationRequest) (res sdk_interface.TeleportationView, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":       "hdrAPINewTeleportation",
		"target":        req.Target,
		"mode":          req.Mode,
		"sourceNetwork": req.SourceNetwork,
		"sourceAddress": req.SourceAddress,
		"sinkNetwork":   req.SinkNetwork,
		"sinkAddress":   req.SinkAddress,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return
	}

	tp, err := mp.NewTeleportation(ctx, target, req.SourceNetwork, req.SourceAddress, req.SinkNetwork, req.SinkAddress, well_known_option.WithMode(req.Mode))
	if err != nil {
		logger.WithError(err).Errorf("failed to new teleportation")
		return
	}

	logger.Infof("new teleportation")

	return ViewTeleportation(tp), nil
}
