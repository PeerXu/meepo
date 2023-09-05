package transport_pipe

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *PipeTransport) setState(s meepo_interface.TransportState) {
	t.state.Store(s)
	if h := t.onTransportStateChange; h != nil {
		h(t)
	}
}
