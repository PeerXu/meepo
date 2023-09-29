package transport_webrtc

import (
	"context"
	"errors"

	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	transport_core "github.com/PeerXu/meepo/pkg/meepo/transport/core"
)

type CloseRequest struct{}

type CloseResponse struct{}

func (t *WebrtcTransport) closeRemoteTransport(ctx context.Context) (err error) {
	var res CloseResponse
	logger := t.GetLogger().WithField("#method", "closeRemoteTransport")
	if err = t.Call(ctx, SYS_METHOD_CLOSE, &CloseRequest{}, &res, well_known_option.WithScope("sys")); err != nil {
		logger.WithError(err).Debugf("failed to close")
		return
	}

	logger.Tracef("remote transport closed")

	return
}

// dont close peer connection now
func (t *WebrtcTransport) onClose(ctx context.Context, _req any) (res any, err error) {
	res = &CloseResponse{}

	logger := t.GetLogger().WithField("#method", "onClose")

	if !t.tryClose() {
		return
	}

	if h := t.BeforeCloseTransportHook; h != nil {
		if er := h(t); er != nil {
			logger.WithError(er).Debugf("before close transport hook failed")
			err = errors.Join(err, er)
		}
	}

	t.readyOnce.Do(func() {
		t.readyErrVal.Store(transport_core.ErrTransportClosed)
		close(t.ready)
	})

	if er := t.closeAllChannels(ctx); er != nil {
		logger.WithError(err).Debugf("failed to close all channels")
		err = errors.Join(err, er)
	}

	if h := t.AfterCloseTransportHook; h != nil {
		h(t)
	}

	logger.Tracef("transport closed")

	return
}
