package listenerer_net

import (
	"context"
	"net"

	listenerer_core "github.com/PeerXu/meepo/pkg/internal/listenerer/core"
	listenerer_interface "github.com/PeerXu/meepo/pkg/internal/listenerer/interface"
)

type listenFunc func(string, string) (net.Listener, error)

func (fn listenFunc) Listen(ctx context.Context, network, address string, opts ...listenerer_interface.ListenOption) (listenerer_interface.Listener, error) {
	lis, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return &listener{lis}, nil
}

var defaultListenFunc listenFunc = net.Listen

func init() {
	listenerer_core.RegisterListenerer("tcp", defaultListenFunc)
}
