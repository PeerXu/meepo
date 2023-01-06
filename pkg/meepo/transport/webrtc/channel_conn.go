package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (c *WebrtcSourceChannel) Conn() meepo_interface.Conn {
	return c.conn
}

func (c *WebrtcSinkChannel) Conn() meepo_interface.Conn {
	if v := c.connVal.Load(); v != nil {
		return v.(meepo_interface.Conn)
	}
	return nil
}
