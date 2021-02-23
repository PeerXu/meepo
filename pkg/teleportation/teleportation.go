package teleportation

import (
	"net"

	"github.com/PeerXu/meepo/pkg/transport"
)

type Portal int

const (
	PortalSource Portal = iota + 1
	PortalSink
)

var (
	PortalStr = []string{
		PortalSource: "source",
		PortalSink:   "sink",
	}
)

func (t Portal) String() string {
	return PortalStr[t]
}

type Teleportation interface {
	Name() string
	Source() net.Addr
	Sink() net.Addr
	Portal() Portal
	Transport() transport.Transport
	DataChannels() []transport.DataChannel
	Close() error
}
