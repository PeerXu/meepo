package listenerer_http

import (
	"net"
	"net/http"
	"sync"
	"time"

	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

type HttpListener struct {
	addr                   net.Addr
	logger                 logging.Logger
	lis                    net.Listener
	server                 *http.Server
	conns                  chan listenerer_interface.Conn
	closeOnce              sync.Once
	connWaitEnabledTimeout time.Duration
	transport              *http.Transport
}
