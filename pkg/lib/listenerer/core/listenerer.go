package listenerer_core

import (
	"context"
	"sync"

	listenerer_interface "github.com/PeerXu/meepo/pkg/lib/listenerer/interface"
)

type (
	Listenerer   = listenerer_interface.Listenerer
	Listener     = listenerer_interface.Listener
	Conn         = listenerer_interface.Conn
	ListenOption = listenerer_interface.ListenOption
)

type listenerer struct{ sync.Map }

func (l *listenerer) Listen(ctx context.Context, network, address string, opts ...ListenOption) (Listener, error) {
	v, ok := l.Load(network)
	if !ok {
		return nil, ErrUnsupportedNetworkFn(network)
	}

	return v.(Listenerer).Listen(ctx, network, address, opts...)
}

var globalListenerer = &listenerer{}

func RegisterListenerer(network string, listenerer Listenerer) {
	globalListenerer.Store(network, listenerer)
}

func GetGlobalListenerer() Listenerer {
	return globalListenerer
}
