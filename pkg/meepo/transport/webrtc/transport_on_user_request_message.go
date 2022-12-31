package transport_webrtc

func (t *WebrtcTransport) onUserRequestMessage(in Message) {
	logger := t.GetLogger().WithField("#method", "onUserRequestMessage").WithFields(t.wrapMessage(in))

	ctx := t.context()
	buf, err := t.onHandle(ctx, in.Method, in.Data)
	if err != nil {
		t.sendErrorResponse(ctx, in, err)
		logger.WithError(err).Debugf("failed to on handle")
		return
	}

	t.sendResponse(ctx, in, buf)

	logger.Tracef("on user request message")
}
