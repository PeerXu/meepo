package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *WebrtcTransport) Addr() meepo_interface.Addr {
	return t.addr
}
