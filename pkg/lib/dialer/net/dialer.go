package dialer_net

import (
	"context"
	"net"

	dialer_core "github.com/PeerXu/meepo/pkg/lib/dialer/core"
	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
)

type dialFunc func(string, string) (net.Conn, error)

func (fn dialFunc) Dial(ctx context.Context, network, address string, opts ...dialer_interface.DialOption) (dialer_interface.Conn, error) {
	return fn(network, address)
}

var defaultDialFunc dialFunc = net.Dial

func init() {
	dialer_core.RegisterDialer("tcp", defaultDialFunc)
}
