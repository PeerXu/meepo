package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (mp *Meepo) NewChannel(ctx context.Context, target Addr, network, address string, opts ...NewChannelOption) (Channel, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method": "NewChannel",
		"target":  target,
		"network": network,
		"address": address,
	})

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		logger.WithError(err).Debugf("failed to get transport")
		return nil, err
	}

	c, err := t.NewChannel(ctx, network, address, opts...)
	if err != nil {
		logger.WithError(err).Debugf("failed to new channel")
		return nil, err
	}

	logger.Tracef("new channel")

	return c, nil
}
