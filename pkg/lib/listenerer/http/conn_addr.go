package listenerer_http

import (
	"net"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
)

func (c *HttpConn) LocalAddr() net.Addr {
	return dialer.NewAddr("tcp", c.request.RemoteAddr)
}

func (c *HttpConn) RemoteAddr() net.Addr {
	return dialer.NewAddr("tcp", c.request.URL.Host)
}
