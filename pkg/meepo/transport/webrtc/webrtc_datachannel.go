//go:build !(js && wasm)

package transport_webrtc

import (
	"github.com/pion/webrtc/v3"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
)

func (t *WebrtcTransport) afterRawSourceChannelCreate(dc *webrtc.DataChannel, c *WebrtcSourceChannel) {
	dc.OnOpen(func() {
		var err error
		if c.conn, err = dc.Detach(); err != nil {
			panic(err)
		}
		c.setState(meepo_interface.CHANNEL_STATE_OPEN)
		c.ready()
	})
}
