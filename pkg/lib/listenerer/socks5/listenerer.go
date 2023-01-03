package listenerer_socks5

import (
	"context"
	"net"

	"github.com/things-go/go-socks5"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	listenerer_core "github.com/PeerXu/meepo/pkg/lib/listenerer/core"
	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

type Socks5Listenerer struct{}

func (l *Socks5Listenerer) Listen(ctx context.Context, network, address string, opts ...listenerer_interface.ListenOption) (listenerer_interface.Listener, error) {
	o := option.Apply(opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	sl := &Socks5Listener{
		addr:   dialer.NewAddr(network, address),
		lis:    lis,
		logger: logger,
		conns:  make(chan *Socks5Conn),
	}
	sl.server = socks5.NewServer(
		socks5.WithConnectHandle(sl.onConnect),
		socks5.WithLogger(logger),
	)
	go sl.Serve(lis) // nolint:errcheck

	return sl, nil
}

func init() {
	listenerer_core.RegisterListenerer("socks5", &Socks5Listenerer{})
}
