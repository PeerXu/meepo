package listenerer_interface

import "net"

type Listener interface {
	Accept() (Conn, error)
	Close() error
	Addr() net.Addr
}
