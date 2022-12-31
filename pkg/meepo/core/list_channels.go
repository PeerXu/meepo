package meepo_core

import (
	"context"
)

func (mp *Meepo) ListChannels(ctx context.Context, opts ...ListChannelsOption) (cm map[Addr][]Channel, err error) {
	logger := mp.GetLogger().WithField("#method", "ListChannels")

	ts, err := mp.ListTransports(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to list transports")
		return nil, err
	}

	cm = make(map[Addr][]Channel)
	for _, t := range ts {
		logger := logger.WithField("addr", t.Addr())
		cs, err := t.ListChannels(ctx, opts...)
		if err != nil {
			logger.WithError(err).Debugf("failed to list channels")
			return nil, err
		}
		cm[t.Addr()] = cs
	}

	logger.Tracef("list channels")

	return cm, nil
}
