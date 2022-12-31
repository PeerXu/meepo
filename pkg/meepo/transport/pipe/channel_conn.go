package transport_pipe

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *PipeChannel) Conn() meepo_interface.Conn {
	return c.conn
}
