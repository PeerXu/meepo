package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) GetChannel(ctx context.Context, target Addr, id uint16) (Channel, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":   "GetChannel",
		"target":    target,
		"channelID": id,
	})

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return nil, err
	}

	c, err := t.GetChannel(ctx, id)
	if err != nil {
		logger.WithError(err).Debugf("failed to get channel")
		return nil, err
	}

	logger.Tracef("get channel")

	return c, nil
}
