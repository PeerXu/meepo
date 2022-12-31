package teleportation_core

import (
	"net"
)

func (tp *teleportation) SinkAddr() net.Addr {
	return tp.sinkAddr
}
