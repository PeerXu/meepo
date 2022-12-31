package transport_webrtc

import (
	"context"

	"github.com/PeerXu/meepo/pkg/internal/well_known_option"
)

type CloseRequest struct{}

type CloseResponse struct{}

func (t *WebrtcTransport) closeRemoteTransport(ctx context.Context) (err error) {
	var res CloseResponse
	logger := t.GetLogger().WithField("#method", "closeRemoteTransport")
	if err = t.Call(ctx, "close", &CloseRequest{}, &res, well_known_option.WithScope("sys")); err != nil {
		logger.WithError(err).Debugf("failed to close")
		return
	}

	logger.Tracef("remote transport closed")

	return
}

func (t *WebrtcTransport) onClose(ctx context.Context, _req any) (res any, err error) {
	res = &CloseResponse{}

	logger := t.GetLogger().WithField("#method", "onClose")

	logger.WithError(t.closeAllChannels(ctx)).Tracef("close all channels")

	if err = t.onCloseCb(t); err != nil {
		logger.WithError(err).Debugf("failed on close callback")
	}

	return
}
