package tracker_transport

import "github.com/PeerXu/meepo/pkg/lib/addr"

func (tk *TransportTracker) Addr() addr.Addr {
	return tk.transport.Addr()
}
