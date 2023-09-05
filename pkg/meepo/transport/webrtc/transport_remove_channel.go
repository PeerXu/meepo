package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *WebrtcTransport) removeChannel(c meepo_interface.Channel) {
	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	t.removeChannelNL(c)
}

func (t *WebrtcTransport) removeChannelNL(c meepo_interface.Channel) {
	delete(t.cs, c.ID())
}
