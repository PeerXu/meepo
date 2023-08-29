package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) hdrAPICloseChannel(ctx context.Context, _req any) (any, error) {
	req := _req.(*sdk_interface.CloseChannelRequest)

	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":   "hdrAPICloseChannel",
		"target":    req.Target,
		"channelID": req.ChannelID,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return nil, err
	}

	c, err := mp.GetChannel(ctx, target, req.ChannelID)
	if err != nil {
		logger.WithError(err).Errorf("failed to get channel")
		return nil, err
	}

	if err = c.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close channel")
		return nil, err
	}

	logger.Infof("channel closed")

	return rpc_core.NO_CONTENT, nil
}
