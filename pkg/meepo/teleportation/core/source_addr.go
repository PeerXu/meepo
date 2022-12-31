package teleportation_core

import "net"

func (tp *teleportation) SourceAddr() net.Addr {
	return tp.sourceAddr
}
