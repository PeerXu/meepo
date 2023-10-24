package listenerer_http

import (
	"net"
	"net/http"
	"sync"

	"github.com/PeerXu/meepo/pkg/lib/logging"
)

type HttpListener struct {
	addr      net.Addr
	logger    logging.Logger
	lis       net.Listener
	server    *http.Server
	conns     chan *HttpConn
	closeOnce sync.Once
}
