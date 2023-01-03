package listenerer_socks5

import (
	"net"
	"strconv"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
)

func (c *Socks5Conn) LocalAddr() net.Addr {
	return c.request.LocalAddr
}

func (c *Socks5Conn) RemoteAddr() net.Addr {
	var host string
	if c.request.DestAddr.FQDN != "" {
		host = c.request.DestAddr.FQDN
	} else {
		host = c.request.DestAddr.IP.String()
	}
	address := net.JoinHostPort(host, strconv.Itoa(c.request.DestAddr.Port))
	return dialer.NewAddr("tcp", address)
}
