package sdk_interface

import (
	"net"

	"github.com/PeerXu/meepo/pkg/lib/addr"
	"github.com/PeerXu/meepo/pkg/lib/version"
)

type SDK interface {
	GetVersion() (*version.V, error)
	Whoami() (addr.Addr, error)
	Ping(target addr.Addr, nonce uint32) (uint32, error)
	Diagnostic() (DiagnosticReportView, error)

	NewTransport(target addr.Addr, opts ...NewTransportOption) (TransportView, error)
	CloseTransport(target addr.Addr) error
	GetTransport(target addr.Addr) (TransportView, error)
	ListTransports() ([]TransportView, error)
	WatchTransports() (<-chan TransportView, <-chan error, func(), error)

	CloseChannel(target addr.Addr, id uint16) error
	GetChannel(target addr.Addr, id uint16) (ChannelView, error)
	ListChannels() ([]ChannelView, error)
	ListChannelsByTarget(target addr.Addr) ([]ChannelView, error)

	NewTeleportation(target addr.Addr, sourceAddr, sinkAddr net.Addr, mode string) (TeleportationView, error)
	CloseTeleportation(id string) error
	GetTeleportation(id string) (TeleportationView, error)
	ListTeleportations() ([]TeleportationView, error)
	Teleport(target addr.Addr, sourceAddr, sinkAddr net.Addr, opts ...TeleportOption) (TeleportationView, error)
}
