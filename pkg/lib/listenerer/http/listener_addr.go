package listenerer_http

import "net"

func (l *HttpListener) Addr() net.Addr {
	return l.addr
}
