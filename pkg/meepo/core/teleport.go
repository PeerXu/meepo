package meepo_core

import (
	"context"
	"errors"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

func (mp *Meepo) Teleport(ctx context.Context, target Addr, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress string, opts ...TeleportOption) (Teleportation, error) {
	logger := mp.GetLogger().WithFields(logging.Fields{
		"#method":       "Teleport",
		"target":        target.String(),
		"sourceNetwork": sourceNetwork,
		"sourceAddress": sourceAddress,
		"sinkNetwork":   sinkNetwork,
		"sinkAddress":   sinkAddress,
	})

	t, err := mp.GetTransport(ctx, target)
	if err != nil {
		if !errors.Is(err, ErrTransportNotFound) {
			logger.WithError(err).Debugf("failed to get transport")
			return nil, err
		}

		t, err = mp.NewTransport(ctx, target, opts...)
		if err != nil {
			logger.WithError(err).Debugf("failed to new transport")
			return nil, err
		}
	}
	if err = t.WaitReady(); err != nil {
		logger.WithError(err).Debugf("failed to wait transport ready")
		return nil, err
	}

	tp, err := mp.NewTeleportation(ctx, target, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress, opts...)
	if err != nil {
		logger.WithError(err).Debugf("failed to new teleportation")
		return nil, err
	}

	logger.Tracef("teleport")

	return tp, nil
}
