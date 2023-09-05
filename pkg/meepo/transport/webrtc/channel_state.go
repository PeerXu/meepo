package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *WebrtcChannel) State() meepo_interface.ChannelState {
	return c.s.Load()
}
