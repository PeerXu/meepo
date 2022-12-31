package listenerer_interface

import (
	"io"
	"net"
)

type Conn interface {
	io.ReadWriteCloser
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
