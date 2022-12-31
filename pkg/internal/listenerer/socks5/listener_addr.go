package listenerer_socks5

import "net"

func (l *Socks5Listener) Addr() net.Addr {
	return l.addr
}
