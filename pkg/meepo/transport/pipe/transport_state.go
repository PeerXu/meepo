package transport_pipe

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *PipeTransport) State() meepo_interface.TransportState {
	return t.state.Load()
}
