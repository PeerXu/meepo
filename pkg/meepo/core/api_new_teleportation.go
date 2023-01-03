package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPINewTeleportation(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.NewTeleportationRequest)
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
		return nil, err
	}

	tp, err := mp.NewTeleportation(ctx, target, req.SourceNetwork, req.SourceAddress, req.SinkNetwork, req.SinkAddress, well_known_option.WithMode(req.Mode))
	if err != nil {
		logger.WithError(err).Errorf("failed to new teleportation")
		return nil, err
	}

	logger.Infof("new teleportation")

	return ViewTeleportation(tp), nil
}
