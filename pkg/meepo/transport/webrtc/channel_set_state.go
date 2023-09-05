package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *WebrtcSourceChannel) setState(s meepo_interface.ChannelState) {
	c.s.Store(s)
	if h := c.onStateChange; h != nil {
		h(c)
	}
}

func (c *WebrtcSinkChannel) setState(s meepo_interface.ChannelState) {
	c.s.Store(s)
	if h := c.onStateChange; h != nil {
		h(c)
	}
}
