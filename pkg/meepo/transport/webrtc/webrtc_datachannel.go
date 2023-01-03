//go:build !(js && wasm)

package transport_webrtc

import "github.com/pion/webrtc/v3"

func (t *WebrtcTransport) afterRawSourceChannelCreate(dc *webrtc.DataChannel, c *WebrtcSourceChannel) {
	dc.OnOpen(func() {
		var err error

		t.csMtx.Lock()
		defer t.csMtx.Unlock()
		t.cs[c.ID()] = c

		if c.conn, err = dc.Detach(); err != nil {
			panic(err)
		}

		c.readyOnce.Do(func() { close(c.ready) })
	})
}
