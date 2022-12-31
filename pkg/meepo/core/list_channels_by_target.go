package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) ListChannelsByTarget(ctx context.Context, target Addr, opts ...ListChannelsOption) ([]Channel, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "ListChannelsByTarget",
		"target":  target.String(),
	})

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return nil, err
	}

	cs, err := t.ListChannels(ctx, opts...)
	if err != nil {
		logger.WithError(err).Debugf("failed to list channels")
		return nil, err
	}

	logger.Tracef("list channels")

	return cs, nil
}
