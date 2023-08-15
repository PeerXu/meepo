package meepo_core

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_eventloop_core "github.com/PeerXu/meepo/pkg/meepo/eventloop/core"
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

	mp.eventloop.Emit(meepo_eventloop_core.NewEvent(EVENT_CHANNEL_ACTION_NEW, nil))

	logger.Tracef("new channel")

	return c, nil
}
