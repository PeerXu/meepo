//go:build js && wasm

package transport_webrtc

import "github.com/pion/webrtc/v3"

func (t *WebrtcTransport) afterRawSourceChannelCreate(dc *webrtc.DataChannel, c *WebrtcSourceChannel) {
	f := func() {
		var err error

		t.csMtx.Lock()
		defer t.csMtx.Unlock()

		t.addChannelNL(c)

		if c.conn, err = dc.Detach(); err != nil {
			panic(err)
		}

		c.ready()
	}
	dc.OnOpen(f)

	// HACK: pion webrtc dont emit callback after data channel has been opened
	if dc.ReadyState() == webrtc.DataChannelStateOpen {
		go f()
	}
}
