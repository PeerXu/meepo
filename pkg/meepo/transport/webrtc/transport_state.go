package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *WebrtcTransport) State() meepo_interface.TransportState {
	return meepo_interface.TransportState(t.pc.ConnectionState().String())
}
