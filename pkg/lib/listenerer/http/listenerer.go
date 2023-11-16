package listenerer_http

import (
	"context"
	"net"
	"net/http"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	listenerer_core "github.com/PeerXu/meepo/pkg/lib/listenerer/core"
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

type HttpListenerer struct{}

func (l *HttpListenerer) Listen(ctx context.Context, network, address string, opts ...listenerer_interface.ListenOption) (listenerer_interface.Listener, error) {
	o := option.ApplyWithDefault(DefaultListenOption(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	connWaitEnabledTimeout, _ := well_known_option.GetConnWaitEnabledTimeoout(o)

	hl := &HttpListener{
		addr:                   dialer.NewAddr(network, address),
		lis:                    lis,
		logger:                 logger,
		conns:                  make(chan listenerer_interface.Conn),
		connWaitEnabledTimeout: connWaitEnabledTimeout,
	}
	hl.transport = &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			p1, p2 := NewHttpGetConn(network, address).Pipe()
			hl.conns <- p1
			if err := p1.WaitEnabled(hl.connWaitEnabledTimeout); err != nil {
				return nil, err
			}
			return p2, nil
		},
	}

	go http.Serve(hl.lis, hl) // nolint:errcheck

	return hl, nil
}

func init() {
	listenerer_core.RegisterListenerer(NAME, &HttpListenerer{})
}
