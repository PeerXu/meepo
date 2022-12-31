package transport_webrtc

import "context"

func (t *WebrtcTransport) sendErrorResponse(ctx context.Context, in Message, err error) {
	logger := t.GetLogger().WithField("#method", "sendErrorResponse").WithFields(t.wrapMessage(in))

	out := t.NewErrorResponse(in, err)
	if er := t.sendMessage(ctx, out); er != nil {
		logger.WithError(er).Debugf("failed to send error response")
		return
	}

	logger.Tracef("send error response")
}
