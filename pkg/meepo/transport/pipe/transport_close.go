package transport_pipe

import (
	"context"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *PipeTransport) Close(ctx context.Context) error {
	logger := t.GetLogger().WithField("#method", "Close")

	if h := t.BeforeCloseTransportHook; h != nil {
		if err := h(t); err != nil {
			logger.WithError(err).Debugf("before close transport hook failed")
			return err
		}
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

	t.setState(meepo_interface.TRANSPORT_STATE_CLOSED)

	if h := t.AfterCloseTransportHook; h != nil {
		h(t)
	}

	logger.Tracef("transport closed")

	return nil
}
