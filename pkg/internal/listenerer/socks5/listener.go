package listenerer_socks5

import (
	"net"
	"sync"

	"github.com/things-go/go-socks5"

	"github.com/PeerXu/meepo/pkg/internal/logging"
)

type Socks5Listener struct {
	addr      net.Addr
	logger    logging.Logger
	lis       net.Listener
	server    *socks5.Server
	conns     chan *Socks5Conn
	closeOnce sync.Once
}
