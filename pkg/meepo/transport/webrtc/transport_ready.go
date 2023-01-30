package transport_webrtc

func (t *WebrtcTransport) Ready() <-chan struct{} {
	return t.ready
}

func (t *WebrtcTransport) isReady() bool {
	select {
	case <-t.Ready():
		return true
	default:
		return false
	}
}
