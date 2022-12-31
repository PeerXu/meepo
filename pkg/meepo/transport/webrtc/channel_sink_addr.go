package transport_webrtc

import "net"

func (c *WebrtcChannel) SinkAddr() net.Addr {
	return c.sinkAddr
}
