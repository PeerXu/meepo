package transport_pipe

import (
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *PipeTransport) Addr() meepo_interface.Addr {
	return t.addr
}
