package transport_webrtc

func (t *WebrtcTransport) onMessage(data []byte) {
	var m Message

	logger := t.GetLogger().WithField("#method", "onMessage")

	err := t.unmarshaler.Unmarshal(data, &m)
	if err != nil {
		logger.WithField("data", string(data)).WithError(err).Debugf("failed to unmarshal data channel message")
		return
	}

	logger = logger.WithFields(t.wrapMessage(m))

	if t.isResponseSession(m.Session) {
		t.onResponseMessage(m)
	} else {
		t.onRequestMessage(m)
	}

	logger.Tracef("on message")
}
