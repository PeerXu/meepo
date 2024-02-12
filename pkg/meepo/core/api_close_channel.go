package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) apiCloseChannel(ctx context.Context, req sdk_interface.CloseChannelRequest) (res rpc_core.EMPTY, err error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":   "hdrAPICloseChannel",
		"target":    req.Target,
		"channelID": req.ChannelID,
	})

	target, err := addr.FromString(req.Target)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse target")
		return
	}

	c, err := mp.GetChannel(ctx, target, req.ChannelID)
	if err != nil {
		logger.WithError(err).Errorf("failed to get channel")
		return
	}

	if err = c.Close(ctx); err != nil {
		logger.WithError(err).Errorf("failed to close channel")
		return
	}

	logger.Infof("channel closed")

	return rpc_core.NO_CONTENT, nil
}
