package transport_webrtc

func (c *WebrtcChannel) IsSource() bool { return false }

func (c *WebrtcSourceChannel) IsSource() bool { return true }
