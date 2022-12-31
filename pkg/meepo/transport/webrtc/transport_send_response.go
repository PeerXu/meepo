package transport_webrtc

import "context"

func (t *WebrtcTransport) sendResponse(ctx context.Context, in Message, data []byte) {
	logger := t.GetLogger().WithField("#method", "sendResponse").WithFields(t.wrapMessage(in))

	out := t.NewResponse(in, data)
	if err := t.sendMessage(ctx, out); err != nil {
		logger.WithError(err).Debugf("failed to send response")
		return
	}

	logger.Tracef("send response")
}
