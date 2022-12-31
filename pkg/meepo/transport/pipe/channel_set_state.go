package transport_pipe

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *PipeChannel) setState(s meepo_interface.ChannelState) {
	c.state.Store(s)
}
