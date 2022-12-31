package transport_webrtc

func (t *WebrtcTransport) onRequestMessage(in Message) {
	logger := t.GetLogger().WithField("#method", "onRequestMessage").WithFields(t.wrapMessage(in))

	switch in.Scope {
	case "sys":
		t.onSystemRequestMessage(in)
	case "usr":
		t.onUserRequestMessage(in)
	default:
		err := ErrUnsupportedScopeFn(in.Scope)
		t.sendErrorResponse(t.context(), in, err)
		logger.WithError(err).Debugf("unsupported scope")
		return
	}

	logger.Tracef("on request message")

}
