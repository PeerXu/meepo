package meepo_socks5

import (
	"net"

	"github.com/things-go/go-socks5"

	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

const (
	ROOT_DOMAIN = ".mpo"
)

type Socks5Server struct {
	logger   logging.Logger
	mp       meepo_interface.Meepo
	server   *socks5.Server
	listener net.Listener
	errors   chan error
	root     string
}

func NewSocks5Server(opts ...NewSocks5ServerOption) (*Socks5Server, error) {
	o := option.ApplyWithDefault(defaultNewSocks5ServerOptions(), opts...)

	logger, err := well_known_option.GetLogger(o)
	if err != nil {
		return nil, err
	}

	mp, err := meepo_interface.GetMeepo(o)
	if err != nil {
		return nil, err
	}

	lis, err := well_known_option.GetListener(o)
	if err != nil {
		return nil, err
	}

	ss := &Socks5Server{
		logger:   logger,
		mp:       mp,
		listener: lis,
		errors:   make(chan error),
		root:     ROOT_DOMAIN,
	}

	ss.server = socks5.NewServer(
		socks5.WithResolver(ss),
		socks5.WithConnectHandle(ss.onConnect),
		socks5.WithLogger(logger),
	)

	return ss, nil
}
