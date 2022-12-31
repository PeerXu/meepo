package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *WebrtcSourceChannel) State() meepo_interface.ChannelState {
	select {
	case <-c.ready:
		return meepo_interface.ChannelState(c.dc.ReadyState().String())
	default:
		return meepo_interface.CHANNEL_STATE_CONNECTING
	}
}

func (c *WebrtcSinkChannel) State() meepo_interface.ChannelState {
	select {
	case <-c.ready:
		return meepo_interface.ChannelState(c.dc.ReadyState().String())
	default:
		return meepo_interface.CHANNEL_STATE_CONNECTING
	}
}
