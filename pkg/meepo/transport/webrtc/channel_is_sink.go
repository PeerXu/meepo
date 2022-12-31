package transport_webrtc

func (tp *WebrtcChannel) IsSink() bool { return false }

func (tp *WebrtcSinkChannel) IsSink() bool { return true }
