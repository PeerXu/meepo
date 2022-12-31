package transport_webrtc

func (t *WebrtcTransport) Ready() <-chan struct{} {
	return t.ready
}
