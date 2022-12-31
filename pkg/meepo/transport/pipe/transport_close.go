package transport_pipe

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *PipeTransport) Close(ctx context.Context) error {
	logger := t.GetLogger().WithField("#method", "Close")

	defer t.setState(meepo_interface.TRANSPORT_STATE_CLOSED)

	if err := t.onClose(t); err != nil {
		logger.WithError(err).Debugf("failed to onClose")
		return err
	}

	cs, err := t.ListChannels(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to list channels")
		return err
	}

	for _, c := range cs {
		if err = c.Close(ctx); err != nil {
			logger.WithError(err).Debugf("failed to close channel")
			return err
		}
	}

	logger.Tracef("transport closed")

	return nil
}
