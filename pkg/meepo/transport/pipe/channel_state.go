package transport_pipe

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *PipeChannel) State() meepo_interface.ChannelState {
	return c.state.Load().(meepo_interface.ChannelState)
}
