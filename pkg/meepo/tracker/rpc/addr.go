package tracker_rpc

import "github.com/PeerXu/meepo/pkg/lib/addr"

func (tk *RPCTracker) Addr() addr.Addr {
	return tk.addr
}
