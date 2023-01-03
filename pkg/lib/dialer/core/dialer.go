package dialer_core

import (
	"context"
	"sync"

	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
)

type (
	Dialer     = dialer_interface.Dialer
	DialOption = dialer_interface.DialOption
	Conn       = dialer_interface.Conn
)

type dialer struct{ sync.Map }

func (d *dialer) Dial(ctx context.Context, network, address string, opts ...DialOption) (Conn, error) {
	v, ok := d.Load(network)
	if !ok {
		return nil, ErrUnsupportedNetworkFn(network)
	}

	return v.(Dialer).Dial(ctx, network, address, opts...)
}

var globalDialer = &dialer{}

func RegisterDialer(network string, dialer Dialer) {
	globalDialer.Store(network, dialer)
}

func GetGlobalDialer() Dialer {
	return globalDialer
}
