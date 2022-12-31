package teleportation_core

import "github.com/PeerXu/meepo/pkg/lib/addr"

func (tp *teleportation) Addr() addr.Addr {
	return tp.addr
}
