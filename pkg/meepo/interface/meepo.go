package meepo_interface

import (
	"context"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	rpc_interface "github.com/PeerXu/meepo/pkg/lib/rpc/interface"
)

type Addr = addr.Addr

type Meepo interface {
	Addr() Addr
	Close(ctx context.Context) error

	Transporter

	Teleportationer
	Teleport(ctx context.Context, transportAddr Addr, sourceNetwork, sourceAddress, sinkNetwork, sinkAddress string, opts ...TeleportOption) (Teleportation, error)

	NewChannel(ctx context.Context, target Addr, network string, address string, opts ...NewChannelOption) (Channel, error)
	ListChannels(ctx context.Context, opts ...ListChannelsOption) (map[Addr][]Channel, error)
	GetChannel(ctx context.Context, target Addr, id uint16) (Channel, error)

	AsTrackerdHandler() rpc_interface.Handler
	AsAPIHandler() rpc_interface.Handler
}
