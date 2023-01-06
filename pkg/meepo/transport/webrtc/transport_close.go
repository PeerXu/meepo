package transport_webrtc

import (
	"context"

	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) Close(ctx context.Context) (err error) {
	logger := t.GetLogger().WithField("#method", "Close")

	if t.closed.Load().(bool) {
		return nil
	}
	t.closed.Swap(true)

	t.readyOnce.Do(func() {
		t.readyErrVal.Store(transport_core.ErrTransportClosed)
		close(t.ready)
	})

	if er := t.closeRemoteTransport(ctx); er != nil {
		err = er
		logger.WithError(err).Debugf("failed to close remote transport")
	}

	if er := t.closeAllChannels(ctx); er != nil {
		err = er
		logger.WithError(err).Debugf("failed to close all channels")
	}

	if er := t.pc.Close(); er != nil {
		err = er
		logger.WithError(err).Debugf("failed to close peer connection")
	}

	if er := t.onCloseCb(t); er != nil {
		err = er
		logger.WithError(err).Debugf("failed on close callback")
	}

	logger.Tracef("transport closed")

	return
}

func (t *WebrtcTransport) closeAllChannels(ctx context.Context) error {
	logger := t.GetLogger().WithField("#method", "closeAllChannels")

	cs, err := t.ListChannels(ctx)
	if err != nil {
		logger.WithError(err).Debugf("failed to list channels")
		return err
	}

	for _, c := range cs {
		if er := c.Close(ctx); er != nil {
			err = er
			logger.WithField("channelID", c.ID()).WithError(er).Debugf("failed to close channel")
		}
	}

	logger.Tracef("close all channels")

	return err
}
