package transport_webrtc

func (t *WebrtcTransport) triggerStateChange() {
	if h := t.onTransportStateChange; h != nil {
		if old, new := t.prevState.Load(), t.State(); old != new && t.prevState.CompareAndSwap(old, new) {
			h(t)
		}
	}
}
