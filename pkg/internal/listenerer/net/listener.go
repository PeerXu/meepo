package listenerer_net

import (
	"net"

	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
)

type listener struct {
	net.Listener
}

func (l *listener) Accept() (listenerer_interface.Conn, error) {
	return l.Listener.Accept()
}
