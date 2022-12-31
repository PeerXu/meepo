package transport_pipe

import "net"

func (c *PipeChannel) SinkAddr() net.Addr {
	return c.sinkAddr
}
