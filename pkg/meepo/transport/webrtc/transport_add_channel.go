package transport_webrtc

import meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"

func (t *WebrtcTransport) addChannel(c meepo_interface.Channel) {
	t.csMtx.Lock()
	defer t.csMtx.Unlock()

	t.addChannelNL(c)
}

func (t *WebrtcTransport) addChannelNL(c meepo_interface.Channel) {
	t.cs[c.ID()] = c
}
