package transport_webrtc

import (
	"context"
	"time"

	"github.com/pion/webrtc/v3"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

func (t *WebrtcTransport) Close(ctx context.Context) (err error) {
	logger := t.GetLogger().WithField("#method", "Close")

	if t.isClosed() {
		return nil
	}
	t.closed.Store(true)

	if h := t.BeforeCloseTransportHook; h != nil {
		if er := h(t); er != nil {
			logger.WithError(er).Debugf("before close transport hook failed")
		}
	}

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

	if er := t.closeAllPeerConnections(); er != nil {
		err = er
		logger.WithError(err).Debugf("failed to close all peer connections")
	}

	if h := t.AfterCloseTransportHook; h != nil {
		h(t)
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

func (t *WebrtcTransport) isClosed() bool {
	return t.closed.Load()
}

func (t *WebrtcTransport) isClosable() bool {
	logger := t.GetLogger().WithField("#method", "isClosable")

	if t.isClosed() {
		logger.Tracef("transport closed")
		return false
	}

	if t.countPeerConnections() > 0 {
		var sessions []string
		t.peerConnections.Range(func(k Session, v *webrtc.PeerConnection) bool {
			sessions = append(sessions, k.String())
			return true
		})
		logger.WithField("sessions", sessions).Tracef("peer connection counts greater than 0")
		return false
	}

	if t.State() == meepo_interface.TRANSPORT_STATE_NEW {
		if t.stat.failedSourceConnections == 0 || t.stat.failedSinkConnections == 0 {
			logger.Tracef("waitting for source/sink connection fails")
			return false
		}
	}

	return true
}

func (t *WebrtcTransport) tryCloseFailedTransport() {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method": "tryCloseFailedTransport",
	})

	select {
	case <-t.Ready():
		logger.Tracef("transport is ready")
	case <-time.After(t.gatherTimeout):
		logger.Debugf("gather timeout")
		go t.Close(t.context())
	}
}
