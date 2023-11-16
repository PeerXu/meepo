package listenerer_http

import (
	"net"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
)

func (c *HttpConnectConn) LocalAddr() net.Addr {
	return dialer.NewAddr("tcp", c.request.RemoteAddr)
}

func (c *HttpConnectConn) RemoteAddr() net.Addr {
	return dialer.NewAddr("tcp", c.request.URL.Host)
}

func (c *HttpGetConn) LocalAddr() net.Addr {
	return dialer.NewAddr("meepo", "meepo")
}

func (c *HttpGetConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}
